package engine

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/gandarez/pong-multiplayer-go/internal/geometry"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ballWidth        = 10
	ballInitialSpeed = 2
)

type ball struct {
	width    float64
	position geometry.Vector
	angle    float64
	speed    float64
	bounces  int
}

func newBall(nextSide string) *ball {
	var angle float64
	if nextSide == player1Name {
		angle = -45 + float64(rand.Intn(91))
	} else {
		angle = 135 + float64(rand.Intn(91))
	}

	return &ball{
		angle: angle,
		width: ballWidth,
		speed: ballInitialSpeed,
		position: geometry.Vector{
			X: (ScreenWidth - ballWidth) / 2,
			Y: (ScreenHeight - ballWidth) / 2,
		},
	}
}

func (b *ball) draw(screen *ebiten.Image) {
	for x := 0.0; x < b.width; x++ {
		for y := 0.0; y < b.width; y++ {
			screen.Set(int(x+b.position.X), int(y+b.position.Y), color.RGBA{200, 200, 200, 255})
		}
	}
}

func (b *ball) bounds() geometry.Rect {
	return geometry.Rect{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  b.width,
		Height: b.width,
	}
}

func (b *ball) update() {
	b.position.X += b.speed * math.Cos(b.angle*math.Pi/180)
	b.position.Y += b.speed * math.Sin(b.angle*math.Pi/180)
}

func (b *ball) bounce(p1Bounds, p2Bounds geometry.Rect) {
	if b.position.Y <= ballWidth || b.position.Y >= ScreenHeight-b.width-ballWidth {
		b.angle *= -1
		b.bounces++

		return
	}

	// left bouncer
	if p1Bounds.Intersects(b.bounds()) {
		b.randomBounce()
		b.bounces++

		return
	}

	// right bouncer
	if p2Bounds.Intersects(b.bounds()) {
		b.randomBounce()
		b.bounces++

		return
	}
}

func (b *ball) randomBounce() {
	b.angle = 180 - b.angle - ballWidth + 20*rand.Float64()
}

// checkGoal checks if the ball is out of the field and returns true if it is.
// It also returns the side that ball went out.
func (b *ball) checkGoal() (bool, geometry.Side) {
	if b.position.X+b.width <= 0 {
		// player 2 scores
		return true, geometry.Left
	}

	if b.position.X >= ScreenWidth {
		// player 1 scores
		return true, geometry.Right
	}

	return false, geometry.Undefined
}
