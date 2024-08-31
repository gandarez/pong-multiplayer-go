package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const fieldWidth = 10

func (g *Game) drawField(screen *ebiten.Image) {
	// draw field limits
	for x := 0; x <= ScreenWidth; x++ {
		for y := 0; y <= fieldWidth; y++ {
			screen.Set(x, y, color.RGBA{200, 200, 200, 255})
			screen.Set(x, y+ScreenHeight-fieldWidth, color.RGBA{200, 200, 200, 255})
		}
	}

	// draw delimiter line
	for squareCount, y := 0, 15; squareCount < 30; squareCount++ {
		for w := 0; w < 7; w++ {
			for h := 0; h < 7; h++ {
				screen.Set((ScreenWidth/2)-5+w, h+y, color.RGBA{200, 200, 200, 255})
			}
		}
		y += 17
	}
}
