package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// TwoPlayersState represents the state of the game when two local players are playing.
type TwoPlayersState struct {
	*BasePlayingState
	ball    ball.Ball
	player1 player.Player
	player2 player.Player
	score1  *score
	score2  *score
}

// NewTwoPlayersState creates a new TwoPlayersState.
func NewTwoPlayersState(game *Game) *TwoPlayersState {
	base := NewBasePlayingState(game, game.menu.Level())

	player1 := player.NewLocal("Player 1", geometry.Left, ScreenWidth, ScreenHeight, fieldBorderWidth)
	player2 := player.NewLocal("Player 2", geometry.Right, ScreenWidth, ScreenHeight, fieldBorderWidth)
	ball := ball.NewLocal(geometry.Right, ScreenWidth, ScreenHeight, base.level)
	score1 := newScore(base.game.font, ScreenWidth/2-100)
	score2 := newScore(base.game.font, ScreenWidth/2+50)

	return &TwoPlayersState{
		BasePlayingState: base,
		ball:             ball,
		player1:          player1,
		player2:          player2,
		score1:           score1,
		score2:           score2,
	}
}

// Update updates the game logic.
func (s *TwoPlayersState) Update() error {
	// Update common elements
	if err := s.BasePlayingState.Update(); err != nil {
		return err
	}

	if s.gamePaused {
		return nil
	}

	// Handle player inputs
	input1 := player.Input{
		Up:   ebiten.IsKeyPressed(ebiten.KeyQ),
		Down: ebiten.IsKeyPressed(ebiten.KeyA),
	}
	input2 := player.Input{
		Up:   ebiten.IsKeyPressed(ebiten.KeyUp),
		Down: ebiten.IsKeyPressed(ebiten.KeyDown),
	}
	s.player1.Update(input1)
	s.player2.Update(input2)

	// Update ball
	s.ball.Update(s.player1.Bounds(), s.player2.Bounds())

	// Check for goals
	if goal, side := s.ball.CheckGoal(); goal {
		s.updateScore(side)
		s.ball = s.ball.Reset(side)
		s.player1.Reset()
		s.player2.Reset()
	}

	// Check for winner
	if s.score1.value == maxScore || s.score2.value == maxScore {
		s.game.ChangeState(NewWinnerState(s.game, s.player1, s.player2, s.score1, s.score2))
	}

	return nil
}

// Draw draws the game elements.
func (s *TwoPlayersState) Draw(screen *ebiten.Image) {
	// Draw common elements
	s.BasePlayingState.Draw(screen)

	// Draw players, ball, and scores
	drawPlayer(s.player1, screen)
	drawPlayer(s.player2, screen)
	drawBall(s.ball, screen)
	s.score1.draw(screen)
	s.score2.draw(screen)
}

// updateScore updates the score based on which side the goal was made.
func (s *TwoPlayersState) updateScore(side geometry.Side) {
	if side == geometry.Left {
		s.score1.value++
	} else {
		s.score2.value++
	}
}
