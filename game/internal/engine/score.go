package engine

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/geometry"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
)

type score struct {
	value    int8
	position geometry.Vector
	textFace *text.GoTextFace
}

func (s *score) draw(screen *ebiten.Image) {
	t := ui.Text{
		Value:    strconv.Itoa(int(s.value)),
		FontFace: s.textFace,
		Position: geometry.Vector{
			X: s.position.X,
			Y: s.position.Y,
		},
		Color: color.RGBA{200, 200, 200, 255},
	}

	t.Draw(screen)
}
