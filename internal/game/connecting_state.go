package game

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
)

// ConnectingState represents the state when the game is connecting to the server.
type ConnectingState struct {
	game            *Game
	connectionError error
}

// NewConnectingState creates a new ConnectingState.
func NewConnectingState(game *Game) *ConnectingState {
	return &ConnectingState{
		game: game,
	}
}

func (s *ConnectingState) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.game.networkClient.Close()
		s.game.ChangeState(NewMainMenuState(s.game))
		return nil
	}

	if s.game.networkClient == nil {
		s.connectToServer()
	}

	// Wait until the network game state channel is initialized
	if s.game.networkGameCh != nil {
		select {
		case gameState := <-s.game.networkGameCh:
			// Change to MultiplayerState
			s.game.ChangeState(NewMultiplayerState(s.game, gameState))
		default:
			// No game state received yet
		}
	}

	return nil
}

// Draw draws the connecting state.
func (s *ConnectingState) Draw(screen *ebiten.Image) {
	ui.DrawSplash(screen, s.game.font, ScreenWidth)
	ui.DrawWaitingConnection(screen, s.game.font, ScreenWidth)
}

// connectToServer connects to the game server.
func (s *ConnectingState) connectToServer() {
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
		if err := s.game.networkClient.ReceiveGameState(s.game.networkGameCh); err != nil {
			slog.Error("Error receiving game state", slog.Any("error", err))
			s.connectionError = err
		}
	}()
}
