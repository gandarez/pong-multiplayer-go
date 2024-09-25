package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ai"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// onePlayerState represents the state of the game when playing against the CPU.
type onePlayerState struct {
	ball    ball.Ball
	player1 player.Player
	player2 player.Player
	score1  *score
	score2  *score
	*baseState
}

// newOnePlayerState creates a new onePlayerState.
func newOnePlayerState(game *Game) *onePlayerState {
	base := newBasePlayingState(game, game.menu.Level())

	player1 := player.NewLocal("Player", geometry.Left, ScreenWidth, ScreenHeight, fieldBorderWidth)
	player2 := player.NewLocal("CPU", geometry.Right, ScreenWidth, ScreenHeight, fieldBorderWidth)
	ball := ball.NewLocal(ScreenWidth, ScreenHeight, base.level)
	score1 := newScore1(base.game.font)
	score2 := newScore2(base.game.font)

	return &onePlayerState{
		baseState: base,
		ball:      ball,
		player1:   player1,
		player2:   player2,
		score1:    score1,
		score2:    score2,
	}
}

// update updates the game logic.
func (s *onePlayerState) update() error {
	// update common elements
	s.baseState.update()

	if s.gamePaused {
		return nil
	}

	// handle player input
	input := player.Input{
		Up:   ebiten.IsKeyPressed(ebiten.KeyUp),
		Down: ebiten.IsKeyPressed(ebiten.KeyDown),
	}
	s.player1.Update(input)

	// update CPU player
	s.player2.SetPosition(ai.GuessBallPosition(
		s.ball.Bounds().Y,
		s.player2.Position().Y,
		s.player2.BouncerHeight(),
		ScreenHeight,
		fieldBorderWidth,
	))

	// update ball
	s.updateBallTrail(s.ball)
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
func (s *onePlayerState) draw(screen *ebiten.Image) {
	// draw common elements
	s.baseState.draw(screen)

	// draw players, ball, and scores
	drawPlayer(s.player1.Position(), s.player1.BouncerWidth(), s.player1.BouncerHeight(), screen)
	drawPlayer(s.player2.Position(), s.player2.BouncerWidth(), s.player2.BouncerHeight(), screen)
	drawBall(screen, s.ball.Position(), s.ball.Width(), s.ballTrail)
	s.score1.draw(screen)
	s.score2.draw(screen)
}

// updateScore updates the score based on which side the goal was made.
func (s *onePlayerState) updateScore(side geometry.Side) {
	if side == geometry.Left {
		s.score2.value++
	} else {
		s.score1.value++
	}
}

func (s *onePlayerState) getBall() ball.Ball {
	return s.ball
}

func (*onePlayerState) canPause() bool {
	return true
}
