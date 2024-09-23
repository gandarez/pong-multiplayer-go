package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// WinnerState represents the state when a player has won the game.
type WinnerState struct {
	game   *Game
	winner player.Player
}

// NewWinnerState creates a new WinnerState.
func NewWinnerState(game *Game, p1, p2 player.Player, s1, s2 *score) *WinnerState {
	var winner player.Player
	if s1.value == maxScore {
		winner = p1
	} else {
		winner = p2
	}
	return &WinnerState{
		game:   game,
		winner: winner,
	}
}

// Update updates the winner state.
func (s *WinnerState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.game.ChangeState(NewMainMenuState(s.game))
	}
	return nil
}

// Draw draws the winner screen.
func (s *WinnerState) Draw(screen *ebiten.Image) {
	textFaceLarge, err := s.game.font.Face("ui", 40)
	if err != nil {
		panic(fmt.Errorf("failed to create winner text face: %w", err))
	}

	winnerText := fmt.Sprintf("%s WON", s.winner.Name())
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
