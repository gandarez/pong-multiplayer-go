package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const multiplayerFindOpponentStr = "Finding opponent.."

// DrawWaitingConnection draws the waiting connection screen.
func DrawWaitingConnection(screen *ebiten.Image, font *font.Font, screenWidth float64) {
	textFace, err := font.Face("ui", 20)
	if err != nil {
		panic(err)
	}

	width, _ := text.Measure(multiplayerFindOpponentStr, textFace, 1)

	uiText := Text{
		Value:    multiplayerFindOpponentStr,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (screenWidth - width) / 2,
			Y: 250.0,
		},
		Color: DefaultColor,
	}

	uiText.Draw(screen)
}
