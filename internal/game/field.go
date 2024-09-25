package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
)

// drawField draws the game field, common to all game modes.
func drawField(screen *ebiten.Image) {
	// draw field limits (top and bottom borders)
	for x := 0.; x <= ScreenWidth; x++ {
		for y := 0.; y <= fieldBorderWidth; y++ {
			screen.Set(int(x), int(y), ui.DefaultColor)
			screen.Set(int(x), int(y+ScreenHeight-fieldBorderWidth), ui.DefaultColor)
		}
	}

	// draw delimiter line (dashed)
	for squareCount, y := 0, 15.; squareCount < 30; squareCount++ {
		for w := 0.; w < 7.; w++ {
			for h := 0.; h < 7.; h++ {
				screen.Set(int(ScreenWidth/2-5+w), int(h+y), ui.DefaultColor)
			}
		}

		y += 17.
	}
}
