package game

import (
	"log/slog"

	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
	"github.com/hajimehoshi/ebiten/v2"
)

type MultiplayerState struct {
	*BasePlayingState
	ball    ball.Ball
	player1 player.Player
	player2 player.Player
	score1  *score
	score2  *score
}

func NewMultiplayerState(game *Game, ready network.ReadyMessage) *MultiplayerState {
	base := NewBasePlayingState(game, game.menu.Level())

	// Assign sides based on the ReadyMessage
	var opponentSide geometry.Side
	if ready.Side == geometry.Right {
		opponentSide = geometry.Left
	} else {
		opponentSide = geometry.Right
	}

	player1 := player.NewNetwork(ready.Name, ready.Side, ScreenWidth, ScreenHeight)
	player2 := player.NewNetwork(ready.OpponentName, opponentSide, ScreenWidth, ScreenHeight)
	ball := ball.NewNetwork()
	score1 := newScore(base.game.font, ScreenWidth/2-100)
	score2 := newScore(base.game.font, ScreenWidth/2+50)

	state := &MultiplayerState{
		BasePlayingState: base,
		ball:             ball,
		player1:          player1,
		player2:          player2,
		score1:           score1,
		score2:           score2,
	}

	// Start listening for game state updates
	go state.game.networkClient.ReceiveGameState(state.game.networkGameCh)

	return state
}

func (s *MultiplayerState) Update() error {
	if err := s.BasePlayingState.Update(); err != nil {
		return err
	}

	if s.gamePaused {
		return nil
	}

	up := ebiten.IsKeyPressed(ebiten.KeyUp)
	down := ebiten.IsKeyPressed(ebiten.KeyDown)

	if up || down {
		if err := s.game.networkClient.SendPlayerInput(network.PlayerInput{
			Up:   up,
			Down: down,
		}); err != nil {
			slog.Error("Error sending player input", slog.Any("error", err))
		}
	}

	gameState := <-s.game.networkGameCh

	s.ball.SetPosition(gameState.Ball.Position)
	s.ball.SetAngle(gameState.Ball.Angle)
	s.ball.SetBounces(gameState.Ball.Bounces)

	s.updatePlayerPositions(gameState)
	s.updateScores(gameState)

	if s.score1.value == maxScore || s.score2.value == maxScore {
		s.game.ChangeState(NewWinnerState(s.game, s.player1, s.player2, s.score1, s.score2))
	}

	return nil
}

func (s *MultiplayerState) updatePlayerPositions(gameState network.GameState) {
	// Directly assign positions
	s.player1.SetPosition(gameState.CurrentPlayer.PositionY)
	s.player2.SetPosition(gameState.OpponentPlayer.PositionY)
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
	s.BasePlayingState.Draw(screen)

	drawPlayer(s.player1, screen)
	drawPlayer(s.player2, screen)
	drawBall(s.ball, screen)
	s.score1.draw(screen)
	s.score2.draw(screen)
}
