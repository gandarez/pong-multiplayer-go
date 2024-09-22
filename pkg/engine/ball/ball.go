package ball

import (
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	initialSpeed = 2
	maxSpeed     = 8
	width        = 10
)

type (
	ball struct {
		position geometry.Vector
		width    float64
	}

	// Ball represents a ball.
	Ball interface {
		Angle() float64
		Bounces() int
		Bounds() geometry.Rect
		CheckGoal() (bool, geometry.Side)
		Position() geometry.Vector
		Reset(nextSide geometry.Side) Ball
		SetPosition(pos geometry.Vector)
		Update(p1Bounds, p2Bounds geometry.Rect)
		Width() float64
	}
)

// Bounds returns the bounds of the ball.
func (b *ball) Bounds() geometry.Rect {
	return geometry.Rect{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  b.width,
		Height: b.width,
	}
}
