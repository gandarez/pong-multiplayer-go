package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
)

// mainMenuState represents the main menu of the game.
type mainMenuState struct {
	game *Game
}

// newMainMenuState creates a new mainMenuState.
func newMainMenuState(game *Game) *mainMenuState {
	return &mainMenuState{
		game: game,
	}
}

// update updates the main menu.
func (s *mainMenuState) update() error {
	s.game.menu.Update()

	if s.game.menu.IsReadyToPlay() {
		switch s.game.menu.GameMode() {
		case menu.Undefined:
			return s.game.exit()
		case menu.OnePlayer:
			s.game.changeState(newOnePlayerState(s.game))
		case menu.TwoPlayers:
			s.game.changeState(newTwoPlayersState(s.game))
		case menu.Multiplayer:
			s.game.changeState(NewConnectingState(s.game))
		case menu.Spectator:
			s.game.changeState(newSpectatorState(s.game))
		}
	}

	return nil
}

// draw draws the main menu.
func (s *mainMenuState) draw(screen *ebiten.Image) {
	ui.DrawSplash(screen, s.game.font, ScreenWidth)
	s.game.menu.Draw(screen)
}

func (*mainMenuState) getBall() ball.Ball {
	panic("not implemented")
}

func (*mainMenuState) canPause() bool {
	return false
}
