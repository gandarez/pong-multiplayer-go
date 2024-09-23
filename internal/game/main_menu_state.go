package game

import (
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

// MainMenuState represents the main menu of the game.
type MainMenuState struct {
	game *Game
}

// NewMainMenuState creates a new MainMenuState.
func NewMainMenuState(game *Game) *MainMenuState {
	return &MainMenuState{
		game: game,
	}
}

// Update updates the main menu.
func (s *MainMenuState) Update() error {
	s.game.menu.Update()

	if s.game.menu.IsReadyToPlay() {
		switch s.game.menu.GameMode() {
		case menu.Undefined:
			s.game.Exit()
		case menu.OnePlayer:
			s.game.ChangeState(NewOnePlayerState(s.game))
		case menu.TwoPlayers:
			s.game.ChangeState(NewTwoPlayersState(s.game))
		case menu.Multiplayer:
			s.game.ChangeState(NewConnectingState(s.game))
		}
	}
	return nil
}

// Draw draws the main menu.
func (s *MainMenuState) Draw(screen *ebiten.Image) {
	ui.DrawSplash(screen, s.game.font, ScreenWidth)
	s.game.menu.Draw(screen)
}
