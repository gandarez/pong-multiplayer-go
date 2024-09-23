package game

import (
	"fmt"
	"log/slog"

	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ConnectingState struct {
	game            *Game
	connectionError error
	readyCh         chan network.ReadyMessage // Channel for ready messages
}

func NewConnectingState(game *Game) *ConnectingState {
	return &ConnectingState{
		game:    game,
		readyCh: make(chan network.ReadyMessage), // Initialize the ready message channel
	}
}

func (s *ConnectingState) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.game.networkClient.Close()
		s.game.ChangeState(NewMainMenuState(s.game))
		return nil
	}

	if s.game.networkClient == nil {
		s.connectToServer() // Attempt to connect to server
	}

	select {
	case readyMessage := <-s.readyCh: // Wait for the ReadyMessage
		slog.Info("Received ReadyMessage", slog.Any("readyMessage", readyMessage))
		s.game.ChangeState(NewMultiplayerState(s.game, readyMessage)) // Transition to MultiplayerState
	default:
		// Continue waiting for ReadyMessage
	}

	return nil
}

func (s *ConnectingState) Draw(screen *ebiten.Image) {
	ui.DrawSplash(screen, s.game.font, ScreenWidth)
	ui.DrawWaitingConnection(screen, s.game.font, ScreenWidth)
}

func (s *ConnectingState) connectToServer() {
	s.game.networkGameCh = make(chan network.GameState)
	s.game.networkClient = network.NewClient(s.game.ctx, s.game.cancel, network.BaseURL)
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
		if err := s.game.networkClient.ReceiveReadyMessage(s.readyCh); err != nil {
			slog.Error("Error receiving ready message", slog.Any("error", err))
			s.connectionError = err
		}
	}()
}
