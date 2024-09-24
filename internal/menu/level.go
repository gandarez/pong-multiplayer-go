package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
)

// LevelSelectionState is the state where the player can select the game level.
type LevelSelectionState struct {
	*baseState
}

var _ State = (*LevelSelectionState)(nil)

// NewLevelSelectionState creates a new LevelSelectionState.
func NewLevelSelectionState(menu *Menu) *LevelSelectionState {
	return &LevelSelectionState{
		baseState: &baseState{
			menu: menu,
			options: []string{
				level.Easy.String(),
				level.Medium.String(),
				level.Hard.String(),
				backStr,
			},
		},
	}
}

// Update updates the state.
func (s *LevelSelectionState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.level = level.Easy
			s.menu.readyToPlay = true
		case 1:
			s.menu.level = level.Medium
			s.menu.readyToPlay = true
		case 2:
			s.menu.level = level.Hard
			s.menu.readyToPlay = true
		case 3:
			s.menu.ChangeState(NewMainMenuState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(NewMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *LevelSelectionState) Draw(screen *ebiten.Image) {
	s.drawOptions(screen)
}

// String returns the state name.
func (*LevelSelectionState) String() string {
	return "LevelSelectionState"
}
