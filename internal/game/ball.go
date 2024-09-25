package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// drawBall draws the ball on the screen.
// bx and by are the ball's position.
// width is the ball's width.
func drawBall(screen *ebiten.Image, position geometry.Vector, width float64, trail []geometry.Vector) {
	trailLength := len(trail)
	// draw the trail
	for i, pos := range trail {
		alpha := uint8((i * (255 / trailLength))) // decrease opacity for older positions

		scaledWidth := width * (1.0 - (float64(i)/float64(trailLength))*(float64(i)/float64(trailLength)))

		drawBallAt(screen, pos, scaledWidth, color.RGBA{100, 100, 100, alpha})
	}

	// Draw the current ball position
	drawBallAt(screen, position, width, ui.DefaultColor)
}

// drawBallAt draws a ball at a specific position.
func drawBallAt(screen *ebiten.Image, position geometry.Vector, width float64, clr color.RGBA) {
	if width <= 1 {
		width = 1
	}

	// Create an image for the ball
	ballImg := ebiten.NewImage(int(width), int(width))
	ballImg.Fill(clr)

	// Draw the ball image at the given position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(position.X, position.Y)

	screen.DrawImage(ballImg, op)
}
