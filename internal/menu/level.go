package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
)

// levelSelectionState is the state where the player can select the game level.
type levelSelectionState struct {
	*baseState
}

var _ state = (*levelSelectionState)(nil)

// newLevelSelectionState creates a new levelSelectionState.
func newLevelSelectionState(menu *Menu) *levelSelectionState {
	return &levelSelectionState{
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
func (s *levelSelectionState) Update() {
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
			s.menu.ChangeState(newMainMenuState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(newMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *levelSelectionState) Draw(screen *ebiten.Image) {
	s.drawOptions(screen)
}

// String returns the state name.
func (*levelSelectionState) String() string {
	return "levelSelectionState"
}
