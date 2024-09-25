package game

import (
	"fmt"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// drawPlayer draws a player on the screen.
func drawPlayer(position geometry.Vector, bouncerWidth, bouncerHeight float64, screen *ebiten.Image) {
	// create bouncer image
	playerImg := ebiten.NewImage(int(bouncerWidth), int(bouncerHeight))
	playerImg.Fill(ui.DefaultColor)

	// translate player to the correct position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(position.X, position.Y)

	screen.DrawImage(playerImg, op)
}

func drawPlayerName(name string, namePosition geometry.Vector, screen *ebiten.Image, font *font.Font) error {
	textface, err := font.Face("ui", 12)
	if err != nil {
		slog.Error("failed to get text face to draw player name", slog.Any("error", err))
		return fmt.Errorf("failed to get text face to draw player name: %w", err)
	}

	t := ui.Text{
		Value:    name,
		FontFace: textface,
		Position: namePosition,
		Color:    color.RGBA{200, 200, 200, 200},
	}
	t.Draw(screen)

	return nil
}
