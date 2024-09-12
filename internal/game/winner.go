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

	winnerFaceSource, err := g.assets.NewTextFaceSource("ui")
	if err != nil {
		// panic is not the best way to handle this error, but Draw() does not return an error
		panic(fmt.Errorf("failed to create winner text face source: %w", err))
	}

	font := &text.GoTextFace{
		Source: winnerFaceSource,
		Size:   40,
	}

	winnerText := fmt.Sprintf("%s won", winner)

	positionX, _ := text.Measure(winnerText, font, 1)

	uiText := ui.Text{
		Value:    fmt.Sprintf("%s won", winner),
		FontFace: font,
		Position: geometry.Vector{
			X: (ScreenWidth - positionX) / 2,
			Y: 200,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)

	font.Size = 30

	positionX, _ = text.Measure("Press Enter to play again", font, 1)

	uiText = ui.Text{
		Value:    "Press Enter to play again",
		FontFace: font,
		Position: geometry.Vector{
			X: (ScreenWidth - positionX) / 2,
			Y: 300,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)
}
