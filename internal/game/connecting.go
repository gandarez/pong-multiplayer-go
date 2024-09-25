package game

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
)

// ConnectingState represents the state when the game is connecting to the server.
type ConnectingState struct {
	game            *Game
	networkReadyCh  chan network.ReadyMessage
	connectionError error
}

// NewConnectingState creates a new ConnectingState.
func NewConnectingState(game *Game) *ConnectingState {
	return &ConnectingState{
		game:           game,
		networkReadyCh: make(chan network.ReadyMessage),
	}
}

func (s *ConnectingState) update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.game.networkClient.Close()
		s.game.menu = menu.New(s.game.font, ScreenWidth, ScreenHeight)
		s.game.changeState(newMainMenuState(s.game))

		return nil
	}

	if s.game.networkClient == nil {
		s.connectToServer()
	}

	if s.connectionError != nil {
		slog.Error("failed to connect to server", slog.Any("error", s.connectionError))

		return fmt.Errorf("failed to connect to server: %w", s.connectionError)
	}

	// wait until the network ready channel is initialized
	for ready := range s.networkReadyCh {
		if ready.Ready {
			s.game.changeState(newMultiplayerState(s.game, ready))
		}
	}

	return nil
}

// draw draws the connecting state.
func (s *ConnectingState) draw(screen *ebiten.Image) {
	ui.DrawSplash(screen, s.game.font, ScreenWidth)
	ui.DrawWaitingConnection(screen, s.game.font, ScreenWidth)
}

// connectToServer connects to the game server.
func (s *ConnectingState) connectToServer() {
	s.game.networkClient = network.NewClient(s.game.ctx, s.game.cancel)
	if err := s.game.networkClient.Connect(); err != nil {
		s.connectionError = fmt.Errorf("failed to connect to server: %w", err)
		return
	}

	if err := s.game.networkClient.SendPlayerInfo(network.GameInfo{
		PlayerName:   s.game.menu.PlayerName(),
		Level:        int(s.game.menu.Level()),
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		MaxScore:     maxScore,
	}); err != nil {
		s.connectionError = fmt.Errorf("failed to send player info: %w", err)
		return
	}

	go func() {
		if err := s.game.networkClient.ReceiveReadyMessage(s.networkReadyCh); err != nil {
			slog.Error("error receiving ready message", slog.Any("error", err))
			s.connectionError = err
		}
	}()
}

func (*ConnectingState) getBall() ball.Ball {
	panic("not implemented")
}

func (*ConnectingState) canPause() bool {
	return false
}
