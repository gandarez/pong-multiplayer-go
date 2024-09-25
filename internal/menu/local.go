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

// LocalModeState is the state where the player can select between one or two players.
type LocalModeState struct {
	*baseState
}

var _ State = (*LocalModeState)(nil)

// NewLocalModeState creates a new LocalModeState.
func NewLocalModeState(menu *Menu) *LocalModeState {
	return &LocalModeState{
		baseState: &baseState{
			menu:    menu,
			options: []string{onePlayerStr, twoPlayersStr, backStr},
		},
	}
}

// Update updates the state.
func (s *LocalModeState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.gameMode = OnePlayer
			s.menu.ChangeState(NewLevelSelectionState(s.menu))
		case 1:
			s.menu.gameMode = TwoPlayers
			s.menu.ChangeState(newTwoPlayersInstructionsState(s.menu))
		case 2:
			s.menu.ChangeState(NewMainMenuState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(NewMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *LocalModeState) Draw(screen *ebiten.Image) {
	s.drawOptions(screen)
}

func (*LocalModeState) String() string {
	return "LocalModeState"
}
