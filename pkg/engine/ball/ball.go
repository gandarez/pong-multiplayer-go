package ball

import (
	"errors"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	initialSpeed = 2
	maxSpeed     = 8
	width        = 10
)

// Kind represents the kind of ball.
type Kind int

const (
	// KindLocal represents a local ball.
	KindLocal Kind = iota
	// KindNetwork represents a network ball.
	KindNetwork
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

// New creates a new ball.
// nextSide is the side that the ball will go when the game starts.
// screenWidth and screenHeight are the dimensions of the screen.
// lvl is the level of the game.
func New(kind Kind, nextSide *geometry.Side, screenWidth, screenHeight *float64, lvl *level.Level) (Ball, error) {
	var ball Ball

	switch kind {
	case KindLocal:
		if nextSide == nil || screenWidth == nil || screenHeight == nil || lvl == nil {
			return nil, errors.New("nextSide, screenWidth, screenHeight and lvl are required for local ball")
		}

		ball = newLocal(*nextSide, *screenWidth, *screenHeight, *lvl)
	case KindNetwork:
		ball = newNetwork()
	}

	return ball, nil
}

// Bounds returns the bounds of the ball.
func (b *ball) Bounds() geometry.Rect {
	return geometry.Rect{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  b.width,
		Height: b.width,
	}
}
