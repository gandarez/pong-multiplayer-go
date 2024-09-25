package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// twoPlayersState represents the state of the game when two local players are playing.
type twoPlayersState struct {
	*baseState
	ball    ball.Ball
	player1 player.Player
	player2 player.Player
	score1  *score
	score2  *score
}

// newTwoPlayersState creates a new twoPlayersState.
func newTwoPlayersState(game *Game) *twoPlayersState {
	base := newBasePlayingState(game, game.menu.Level())

	player1 := player.NewLocal("Player 1", geometry.Left, ScreenWidth, ScreenHeight, fieldBorderWidth)
	player2 := player.NewLocal("Player 2", geometry.Right, ScreenWidth, ScreenHeight, fieldBorderWidth)
	ball := ball.NewLocal(ScreenWidth, ScreenHeight, base.level)
	score1 := newScore1(base.game.font)
	score2 := newScore2(base.game.font)

	return &twoPlayersState{
		baseState: base,
		ball:      ball,
		player1:   player1,
		player2:   player2,
		score1:    score1,
		score2:    score2,
	}
}

// update updates the game logic.
func (s *twoPlayersState) update() error {
	// update common elements
	s.baseState.update()

	if s.gamePaused {
		return nil
	}

	// handle player inputs
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

	// update ball
	s.ball.Update(s.player1.Bounds(), s.player2.Bounds())

	// check for goals
	if goal, side := s.ball.CheckGoal(); goal {
		s.updateScore(side)
		s.ball = s.ball.Reset()
		s.player1.Reset()
		s.player2.Reset()
	}

	// check for winner
	if s.score1.value == maxScore || s.score2.value == maxScore {
		winner := s.player1
		if s.score2.value == maxScore {
			winner = s.player2
		}

		s.game.changeState(newWinnerState(s.game, winner.Name(), s))
	}

	return nil
}

// draw draws the game elements.
func (s *twoPlayersState) draw(screen *ebiten.Image) {
	// draw common elements
	s.baseState.draw(screen)

	// draw players, ball, and scores
	drawPlayer(s.player1.Position(), s.player1.BouncerWidth(), s.player1.BouncerHeight(), screen)
	drawPlayer(s.player2.Position(), s.player2.BouncerWidth(), s.player2.BouncerHeight(), screen)
	drawBall(s.ball.Position(), s.ball.Width(), screen)
	s.score1.draw(screen)
	s.score2.draw(screen)
}

// updateScore updates the score based on which side the goal was made.
func (s *twoPlayersState) updateScore(side geometry.Side) {
	if side == geometry.Left {
		s.score2.value++
	} else {
		s.score1.value++
	}
}

func (s *twoPlayersState) getBall() ball.Ball {
	return s.ball
}

func (*twoPlayersState) canPause() bool {
	return true
}
