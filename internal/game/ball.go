package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// drawBall draws the ball on the screen.
// bx and by are the ball's position.
// width is the ball's width.
func drawBall(position geometry.Vector, width float64, screen *ebiten.Image) {
	for x := 0.0; x < width; x++ {
		for y := 0.0; y < width; y++ {
			screen.Set(int(position.X+x), int(position.Y+y), color.RGBA{200, 200, 200, 255})
		}
	}
}
