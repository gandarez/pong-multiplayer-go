package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	localModeStr    = "Local Mode"
	multiplayerStr  = "Multiplayer"
	instructionsStr = "Instructions"
)

// MainMenuState is the state where the player can select between local mode, multiplayer or see the instructions.
type MainMenuState struct {
	*baseState
}

var _ State = (*MainMenuState)(nil)

// NewMainMenuState creates a new MainMenuState.
func NewMainMenuState(menu *Menu) *MainMenuState {
	return &MainMenuState{
		baseState: &baseState{
			menu:    menu,
			options: []string{localModeStr, multiplayerStr, instructionsStr},
		},
	}
}

// Update updates the state.
func (s *MainMenuState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.ChangeState(NewLocalModeState(s.menu))
		case 1:
			s.menu.ChangeState(NewInputNameState(s.menu))
		case 2:
			s.menu.ChangeState(NewInstructionsState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.gameMode = Undefined
		s.menu.readyToPlay = true
	}
}

// Draw draws the state.
func (s *MainMenuState) Draw(screen Screen) {
	s.drawOptions(screen)
}

// String returns the state name.
func (*MainMenuState) String() string {
	return "MainMenuState"
}
