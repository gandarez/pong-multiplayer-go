package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const multiplayerFindOpponentStr = "Finding opponent.."

func DrawWaitingConnection(screen *ebiten.Image, assets *assets.Assets, screenWidth float64) {
	textFaceSource, err := assets.NewTextFaceSource("ui")
	if err != nil {
		panic(err)
		// return nil, fmt.Errorf("failed to create main title text face source: %w", err)
	}

	textFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   20,
	}

	positionX, _ := text.Measure(multiplayerFindOpponentStr, textFace, 1)

	uiText := Text{
		Value:    multiplayerFindOpponentStr,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (screenWidth - positionX) / 2,
			Y: 250.0,
		},
		Color: DefaultColor,
	}

	uiText.Draw(screen)
}
