// internal/game/multiplayer_state.go

package game

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
)

type MultiplayerState struct {
	*BasePlayingState
	ball    ball.Ball
	player1 player.Player
	player2 player.Player
	score1  *score
	score2  *score
}

func NewMultiplayerState(game *Game, gameState network.GameState) *MultiplayerState {
	base := NewBasePlayingState(game, game.menu.Level())

	// Initialize players with names from gameState
	player1 := player.NewNetwork(gameState.CurrentPlayer.Name, gameState.CurrentPlayer.Side, ScreenWidth, ScreenHeight)
	player2 := player.NewNetwork(gameState.OpponentPlayer.Name, gameState.OpponentPlayer.Side, ScreenWidth, ScreenHeight)
	ball := ball.NewNetwork()
	score1 := newScore(base.game.font, ScreenWidth/2-100)
	score2 := newScore(base.game.font, ScreenWidth/2+50)

	return &MultiplayerState{
		BasePlayingState: base,
		ball:             ball,
		player1:          player1,
		player2:          player2,
		score1:           score1,
		score2:           score2,
	}
}

func (s *MultiplayerState) Update() error {
	// Update common elements
	if err := s.BasePlayingState.Update(); err != nil {
		return err
	}

	if s.gamePaused {
		return nil
	}

	up := ebiten.IsKeyPressed(ebiten.KeyUp)
	down := ebiten.IsKeyPressed(ebiten.KeyDown)

	if up || down {
		// Send input to server
		if err := s.game.networkClient.SendPlayerInput(network.PlayerInput{
			Up:   up,
			Down: down,
		}); err != nil {
			slog.Error("Error sending player input", slog.Any("error", err))
		}
	}

	// Receive game state from server and update local game state
	select {
	case gameState := <-s.game.networkGameCh:
		// Update ball and players positions
		s.ball.SetPosition(gameState.Ball.Position)
		s.ball.SetAngle(gameState.Ball.Angle)
		s.ball.SetBounces(gameState.Ball.Bounces)

		// Update player positions and scores
		s.updatePlayerPositions(gameState)
		s.updateScores(gameState)

		// Check for winner
		if s.score1.value == maxScore || s.score2.value == maxScore {
			s.game.ChangeState(NewWinnerState(s.game, s.player1, s.player2, s.score1, s.score2))
		}
	default:
		// No new game state received
	}

	return nil
}

func (s *MultiplayerState) updatePlayerPositions(gameState network.GameState) {
	if s.player1.Side() == gameState.CurrentPlayer.Side {
		s.player1.SetPosition(gameState.CurrentPlayer.PositionY)
		s.player2.SetPosition(gameState.OpponentPlayer.PositionY)
	} else {
		s.player1.SetPosition(gameState.OpponentPlayer.PositionY)
		s.player2.SetPosition(gameState.CurrentPlayer.PositionY)
	}
}

func (s *MultiplayerState) updateScores(gameState network.GameState) {
	if s.player1.Side() == gameState.CurrentPlayer.Side {
		s.score1.value = gameState.CurrentPlayer.Score
		s.score2.value = gameState.OpponentPlayer.Score
	} else {
		s.score1.value = gameState.OpponentPlayer.Score
		s.score2.value = gameState.CurrentPlayer.Score
	}
}

func (s *MultiplayerState) Draw(screen *ebiten.Image) {
	// Draw common elements
	s.BasePlayingState.Draw(screen)

	// Draw players, ball, and scores
	drawPlayer(s.player1, screen)
	drawPlayer(s.player2, screen)
	drawBall(s.ball, screen)
	s.score1.draw(screen)
	s.score2.draw(screen)
}
