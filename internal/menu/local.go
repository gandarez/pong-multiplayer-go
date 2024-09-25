package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	onePlayerStr  = "One Player"
	twoPlayersStr = "Two Players"
	backStr       = "Back"
)

// localModeState is the state where the player can select between one or two players.
type localModeState struct {
	*baseState
}

var _ state = (*localModeState)(nil)

// newLocalModeState creates a new localModeState.
func newLocalModeState(menu *Menu) *localModeState {
	return &localModeState{
		baseState: &baseState{
			menu:    menu,
			options: []string{onePlayerStr, twoPlayersStr, backStr},
		},
	}
}

// Update updates the state.
func (s *localModeState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.gameMode = OnePlayer
			s.menu.ChangeState(newLevelSelectionState(s.menu))
		case 1:
			s.menu.gameMode = TwoPlayers
			s.menu.ChangeState(newTwoPlayersInstructionsState(s.menu))
		case 2:
			s.menu.ChangeState(newMainMenuState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(newMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *localModeState) Draw(screen *ebiten.Image) {
	s.drawOptions(screen)
}

func (*localModeState) String() string {
	return "localModeState"
}
