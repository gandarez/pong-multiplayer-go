package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	localModeStr    = "Local Mode"
	multiplayerStr  = "Multiplayer"
	spectateStr     = "Watch"
	instructionsStr = "Instructions"
)

// mainMenuState is the state where the player can select between local mode, multiplayer or see the instructions.
type mainMenuState struct {
	*baseState
}

var _ state = (*mainMenuState)(nil)

// newMainMenuState creates a new mainMenuState.
func newMainMenuState(menu *Menu) *mainMenuState {
	return &mainMenuState{
		baseState: &baseState{
			menu:    menu,
			options: []string{localModeStr, multiplayerStr, spectateStr, instructionsStr},
		},
	}
}

// Update updates the state.
func (s *mainMenuState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.ChangeState(newLocalModeState(s.menu))
		case 1:
			s.menu.ChangeState(newInputNameState(s.menu))
		case 2:
			s.menu.ChangeState(newSpectateSessionsState(s.menu))
		case 3:
			s.menu.ChangeState(newInstructionsState(s.menu))
		}
	}

	// return to quit the game
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.gameMode = Undefined
		s.menu.readyToPlay = true
	}
}

// Draw draws the state.
func (s *mainMenuState) Draw(screen *ebiten.Image) {
	s.drawOptions(screen)
}

// String returns the state name.
func (*mainMenuState) String() string {
	return "mainMenuState"
}
