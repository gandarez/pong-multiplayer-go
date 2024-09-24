package game

import (
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// drawPlayer draws a player on the screen.
func drawPlayer(position geometry.Vector, bouncerWidth, bouncerHeight float64, screen *ebiten.Image) {
	for x := 0.0; x < bouncerWidth; x++ {
		for y := 0.0; y < bouncerHeight; y++ {
			screen.Set(int(position.X+x), int(position.Y+y), ui.DefaultColor)
		}
	}
}

func drawPlayerName(name string, namePosition geometry.Vector, screen *ebiten.Image, font *font.Font) {
	textface, err := font.Face("ui", 12)
	if err != nil {
		slog.Error("failed to get text face to draw player name", slog.Any("error", err))
		return
	}

	t := ui.Text{
		Value:    name,
		FontFace: textface,
		Position: namePosition,
		Color:    color.RGBA{200, 200, 200, 200},
	}
	t.Draw(screen)
}
