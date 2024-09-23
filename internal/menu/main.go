package menu

import (
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	localModeStr    = "Local Mode"
	multiplayerStr  = "Multiplayer"
	instructionsStr = "Instructions"
)

type MainMenuState struct {
	*BaseState
}

func NewMainMenuState(menu *Menu) *MainMenuState {
	return &MainMenuState{
		BaseState: &BaseState{
			menu:    menu,
			options: []string{localModeStr, multiplayerStr, instructionsStr},
		},
	}
}

func (s *MainMenuState) Update() {
	s.navigateOptions(len(s.options))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.selectedOption {
		case 0:
			s.menu.ChangeState(NewLocalModeState(s.menu))
		case 1:
			if s.menu.playerName != "" {
				s.menu.gameMode = Multiplayer
				s.menu.level = level.Medium
				s.menu.readyToPlay = true
			} else {
				s.menu.ChangeState(NewInputNameState(s.menu))
			}
		case 2:
			s.menu.ChangeState(NewInstructionsState(s.menu))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.gameMode = Undefined
		s.menu.readyToPlay = true
	}
}

func (s *MainMenuState) Draw(screen Screen) {
	s.drawOptions(screen)
}
