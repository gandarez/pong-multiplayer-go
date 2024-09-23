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

type LocalModeState struct {
	*BaseState
}

func NewLocalModeState(menu *Menu) *LocalModeState {
	return &LocalModeState{
		BaseState: &BaseState{
			menu:    menu,
			options: []string{onePlayerStr, twoPlayersStr, backStr},
		},
	}
}

func (s *LocalModeState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.gameMode = OnePlayer
			s.menu.ChangeState(NewLevelSelectionState(s.menu))
		case 1:
			s.menu.gameMode = TwoPlayers
			s.menu.ChangeState(NewLevelSelectionState(s.menu))
		case 2:
			s.menu.ChangeState(NewMainMenuState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(NewMainMenuState(s.menu))
	}
}

func (s *LocalModeState) Draw(screen Screen) {
	s.drawOptions(screen)
}
