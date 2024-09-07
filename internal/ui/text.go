package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// Text represents a text to be drawn on the screen.
type Text struct {
	Value    string
	Position geometry.Vector
	Color    color.RGBA
	FontFace text.Face
}

// Draw draws the text on the screen.
func (t *Text) Draw(screen *ebiten.Image) {
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(t.Position.X, t.Position.Y)
	opts.ColorScale.ScaleWithColor(t.Color)

	text.Draw(screen, t.Value, t.FontFace, opts)
}
