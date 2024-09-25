package game

import (
	"context"
	"fmt"
	"log/slog"

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
	game      *Game
	winner    string
	prevState state
}

// newWinnerState creates a new winnerState.
func newWinnerState(game *Game, winner string, prevState state) *winnerState {
	return &winnerState{
		game:      game,
		winner:    winner,
		prevState: prevState,
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
	// draw previous state
	s.prevState.draw(screen)

	// overlay a semi-transparent layer
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(ui.TransparentBlack)
	screen.DrawImage(overlay, nil)

	textFaceSmall, err := s.game.font.Face("ui", 30)
	if err != nil {
		slog.Error("failed to create winner text face", slog.Any("error", err))
		panic(err)
	}

	if err := s.drawWinner(screen); err != nil {
		slog.Error("failed to draw winner", slog.Any("error", err))
		panic(err)
	}

	instructionText := "Press Enter to play again"
	textWidth, _ := text.Measure(instructionText, textFaceSmall, 1)

	uiText := ui.Text{
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

func (s *winnerState) drawWinner(screen *ebiten.Image) error {
	textFaceLarge, err := s.game.font.Face("ui", 40)
	if err != nil {
		return fmt.Errorf("failed to create winner text face: %w", err)
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

	return nil
}

func (*winnerState) getBall() ball.Ball {
	panic("not implemented")
}

func (*winnerState) canPause() bool {
	return false
}
