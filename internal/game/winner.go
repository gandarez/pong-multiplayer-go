package game

import (
	"context"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// winnerState represents the state when a player has won the game.
type winnerState struct {
	game   *Game
	winner string
}

// newWinnerState creates a new winnerState.
func newWinnerState(game *Game, winner string) *winnerState {
	return &winnerState{
		game:   game,
		winner: winner,
	}
}

// update updates the winner state.
func (s *winnerState) update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.game.menu = menu.New(s.game.font, ScreenWidth, ScreenHeight)
		s.game.networkClient = nil

		ctx, cancel := context.WithCancel(context.Background())
		s.game.ctx = ctx
		s.game.cancel = cancel

		s.game.changeState(newMainMenuState(s.game))
	}

	return nil
}

// draw draws the winner screen.
func (s *winnerState) draw(screen *ebiten.Image) {
	textFaceLarge, err := s.game.font.Face("ui", 40)
	if err != nil {
		panic(fmt.Errorf("failed to create winner text face: %w", err))
	}

	winnerText := fmt.Sprintf("%s WON", s.winner)
	textWidth, _ := text.Measure(winnerText, textFaceLarge, 1)

	uiText := ui.Text{
		Value:    winnerText,
		FontFace: textFaceLarge,
		Position: geometry.Vector{
			X: (ScreenWidth - textWidth) / 2,
			Y: 200,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)

	textFaceSmall, err := s.game.font.Face("ui", 30)
	if err != nil {
		panic(fmt.Errorf("failed to create small text face: %w", err))
	}

	instructionText := "Press Enter to play again"
	textWidth, _ = text.Measure(instructionText, textFaceSmall, 1)

	uiText = ui.Text{
		Value:    instructionText,
		FontFace: textFaceSmall,
		Position: geometry.Vector{
			X: (ScreenWidth - textWidth) / 2,
			Y: 300,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)
}

func (*winnerState) getBall() ball.Ball {
	panic("not implemented")
}

func (*winnerState) canPause() bool {
	return false
}
