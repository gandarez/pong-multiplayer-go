package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	metric "github.com/gandarez/pong-multiplayer-go/internal/stat"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
)

// BasePlayingState contains common logic for playing states.
type BasePlayingState struct {
	game       *Game
	level      level.Level
	pauseMenu  *PauseMenu
	metric     *metric.Metric
	gamePaused bool
}

// NewBasePlayingState creates a new BasePlayingState.
func NewBasePlayingState(game *Game, lvl level.Level) *BasePlayingState {
	pauseMenu := NewPauseMenu(game.font, ScreenWidth)
	metric, err := metric.New(game.font, ScreenWidth)
	if err != nil {
		panic(err)
	}

	return &BasePlayingState{
		game:      game,
		level:     lvl,
		pauseMenu: pauseMenu,
		metric:    metric,
	}
}

// Update handles common update logic, including pause menu.
func (s *BasePlayingState) Update() error {
	// Handle pause menu
	if s.pauseMenu.isOpen {
		s.pauseMenu.Update()
		if s.pauseMenu.ShouldExit {
			s.game.ChangeState(NewMainMenuState(s.game))
			return nil
		}
		if s.pauseMenu.ShouldResume {
			s.gamePaused = false
			s.pauseMenu.ShouldResume = false
		}
		return nil
	}

	// Check for pause input
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.gamePaused = true
		s.pauseMenu.Open()
		return nil
	}

	// Common update logic
	return nil
}

// Draw handles common drawing logic, including pause menu.
func (s *BasePlayingState) Draw(screen *ebiten.Image) {
	// Draw the field
	s.game.drawField(screen)

	// Draw pause menu if open
	if s.gamePaused && s.pauseMenu.isOpen {
		s.pauseMenu.Draw(screen)
	}
}
