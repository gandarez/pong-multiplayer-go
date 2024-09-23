package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	engineball "github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
)

// drawBall draws the ball on the screen.
func drawBall(b engineball.Ball, screen *ebiten.Image) {
	for x := 0.0; x < b.Width(); x++ {
		for y := 0.0; y < b.Width(); y++ {
			screen.Set(int(b.Position().X+x), int(b.Position().Y+y), color.RGBA{200, 200, 200, 255})
		}
	}
}
