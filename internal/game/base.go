package game

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/stat"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
)

// baseState contains common logic for playing states.
type baseState struct {
	game              *Game
	level             level.Level
	pauseMenu         *pauseMenu
	metric            *stat.Metric
	gamePaused        bool
	showMetric        bool
	pingCurrentPlayer int
	pingOpponent      int
}

// newBasePlayingState creates a new baseState.
func newBasePlayingState(game *Game, lvl level.Level) *baseState {
	pauseMenu := newPauseMenu(game.font, ScreenWidth)

	metric, err := stat.New(game.font, ScreenWidth)
	if err != nil {
		slog.Error("failed to create metric", slog.Any("error", err))
	}

	return &baseState{
		game:      game,
		level:     lvl,
		pauseMenu: pauseMenu,
		metric:    metric,
	}
}

// update handles common update logic, including pause menu.
func (s *baseState) update() {
	// handle pause menu
	if s.pauseMenu.isShown {
		s.pauseMenu.update()

		if s.pauseMenu.ShouldExit {
			// force reset the menu
			s.game.menu = menu.New(s.game.font, ScreenWidth)
			s.game.changeState(newMainMenuState(s.game))

			return
		}

		if s.pauseMenu.ShouldResume {
			s.gamePaused = false
			s.pauseMenu.ShouldResume = false
		}

		return
	}

	// if TAB is pressed, then show/hide metrics
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		s.showMetric = !s.showMetric
	}

	// check for pause input
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && s.game.currentState.canPause() {
		s.gamePaused = true
		s.pauseMenu.show()

		return
	}
}

// draw handles common drawing logic, including pause menu.
func (s *baseState) draw(screen *ebiten.Image) {
	// draw the field
	drawField(screen)

	// draw metric if enabled
	s.tryDrawMetric(screen)

	// draw pause menu if open
	if s.gamePaused && s.pauseMenu.isShown {
		s.pauseMenu.draw(screen)
	}
}

func (s *baseState) tryDrawMetric(screen *ebiten.Image) {
	if !s.showMetric {
		return
	}

	if s.metric == nil {
		slog.Warn("metric is nil")

		return
	}

	s.metric.Draw(
		screen,
		s.game.currentState.getBall().Bounces(),
		s.game.currentState.getBall().Angle(),
		s.level,
	)
}
