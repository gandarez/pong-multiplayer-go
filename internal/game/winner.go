package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

func (g *Game) drawWinner(screen *ebiten.Image) {
	var winner string
	if g.score1.value == maxScore {
		winner = g.player1.Name()
	} else {
		winner = g.player2.Name()
	}

	textFace, err := g.font.Face("ui", 40)
	if err != nil {
		// panic is not the best way to handle this error, but Draw() does not return an error
		panic(fmt.Errorf("failed to create winner text face: %w", err))
	}

	winnerText := fmt.Sprintf("%s WON", winner)

	width, _ := text.Measure(winnerText, textFace, 1)

	uiText := ui.Text{
		Value:    winnerText,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (ScreenWidth - width) / 2,
			Y: 200,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)

	textFace.Size = 30

	width, _ = text.Measure("Press Enter to play again", textFace, 1)

	uiText = ui.Text{
		Value:    "Press Enter to play again",
		FontFace: textFace,
		Position: geometry.Vector{
			X: (ScreenWidth - width) / 2,
			Y: 300,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)
}
