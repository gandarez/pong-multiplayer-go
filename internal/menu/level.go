package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
)

type LevelSelectionState struct {
	*BaseState
}

func NewLevelSelectionState(menu *Menu) *LevelSelectionState {
	return &LevelSelectionState{
		BaseState: &BaseState{
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

func (s *LevelSelectionState) Draw(screen Screen) {
	s.drawOptions(screen)
}
