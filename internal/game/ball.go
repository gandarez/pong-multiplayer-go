package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	engineball "github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
)

type ball struct {
	*engineball.Ball
}

func (b ball) draw(screen *ebiten.Image) {
	for x := 0.0; x < b.Width(); x++ {
		for y := 0.0; y < b.Width(); y++ {
			screen.Set(int(x+b.Position().X), int(y+b.Position().Y), color.RGBA{200, 200, 200, 255})
		}
	}
}
