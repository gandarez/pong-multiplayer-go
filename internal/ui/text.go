package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// nolint:revive,gochecknoglobals
var (
	// DefaultColor is the default color of the text.
	DefaultColor = color.RGBA{200, 200, 200, 255}
	// HighlightColor is used to highlight selected options.
	HighlightColor = color.RGBA{255, 255, 0, 255}
	// TransparentBlack is a semi-transparent black color for overlays.
	TransparentBlack = color.RGBA{0, 0, 0, 128}
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
