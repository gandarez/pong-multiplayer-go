package ui

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const titleStr = "PONGO"

// DrawSplash draws the splash screen.
func DrawSplash(screen *ebiten.Image, font *font.Font, screenWidth float64) {
	textFace, err := font.Face("ui", 80)
	if err != nil {
		slog.Error("failed to create main title text face", slog.Any("error", err))

		return
	}

	width, _ := text.Measure(titleStr, textFace, 1)

	uiText := Text{
		Value:    titleStr,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (screenWidth - width) / 2,
			Y: 80,
		},
		Color: DefaultColor,
	}

	uiText.Draw(screen)
}
