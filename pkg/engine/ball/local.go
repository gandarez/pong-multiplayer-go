package ball

import (
	"math"
	"math/rand"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// Local represents the ball in a local game.
type Local struct {
	level        level.Level
	screenHeight float64
	screenWidth  float64
	speed        float64
	*ball
}

// NewLocal creates a new ball to play locally.
// nextSide is the side that the ball will go when the game starts.
// screenWidth and screenHeight are the dimensions of the screen.
// lvl is the level of the game.
func NewLocal(nextSide geometry.Side, screenWidth, screenHeight float64, lvl level.Level) *Local {
	return &Local{
		level:        lvl,
		screenHeight: screenHeight,
		screenWidth:  screenWidth,
		speed:        initialSpeed,
		ball: &ball{
			angle:   calcInitialAngle(nextSide),
			bounces: 0,
			position: geometry.Vector{
				X: (screenWidth - width) / 2,
				Y: (screenHeight - width) / 2,
			},
			width: width,
		},
	}
}

// Angle returns the angle of the ball.
func (b *Local) Angle() float64 {
	return b.angle
}

// Bounces returns the number of bounces of the ball.
func (b *Local) Bounces() int {
	return b.bounces
}

// Bounds returns the bounds of the ball.
func (b *Local) Bounds() geometry.Rect {
	return b.ball.Bounds()
}

// CheckGoal checks if the ball has scored a goal.
func (b *Local) CheckGoal() (bool, geometry.Side) {
	if b.position.X+b.width <= 0 {
		// player 2 scores (right player)
		return true, geometry.Left
	}

	if b.position.X >= b.screenWidth {
		// player 1 scores (left player)
		return true, geometry.Right
	}

	return false, geometry.Undefined
}

// Position returns the position of the ball.
func (b *Local) Position() geometry.Vector {
	return b.position
}

// Reset resets the position of the ball.
func (b *Local) Reset(nextSide geometry.Side) Ball {
	return NewLocal(nextSide, b.screenWidth, b.screenHeight, b.level)
}

// SetAngle will panic because it is not implemented.
func (*Local) SetAngle(_ float64) {
	panic("not implemented")
}

// SetBounces will panic because it is not implemented.
func (*Local) SetBounces(_ int) {
	panic("not implemented")
}

// SetPosition sets the position of the ball.
func (b *Local) SetPosition(pos geometry.Vector) {
	b.position = pos
}

// Update updates the position of the ball.
func (b *Local) Update(p1Bounds, p2Bounds geometry.Rect) {
	b.position.X += b.speed * math.Cos(b.angle*math.Pi/180)
	b.position.Y += b.speed * math.Sin(b.angle*math.Pi/180)

	b.bounce(p1Bounds, p2Bounds)
}

// Width returns the width of the ball.
func (b *Local) Width() float64 {
	return b.width
}

func calcInitialAngle(nextSide geometry.Side) float64 {
	if nextSide == geometry.Left {
		return -45 + float64(rand.Intn(91)) // nolint:gosec
	}

	return 135 + float64(rand.Intn(91)) // nolint:gosec
}

// bounce checks if the ball is bouncing on the walls or the players and changes the angle of the ball.
func (b *Local) bounce(p1Bounds, p2Bounds geometry.Rect) {
	// top wall or bottom wall
	if b.position.Y <= width || b.position.Y >= b.screenHeight-b.width-width {
		b.angle *= -1
		b.bounces++

		b.increaseSpeed()
	}

	// left bouncer or right bouncer
	if p1Bounds.Intersects(b.Bounds()) || p2Bounds.Intersects(b.Bounds()) {
		b.bounces++

		b.randomBounce()
		b.increaseSpeed()
	}
}

// randomBounce changes the angle of the ball randomly.
func (b *Local) randomBounce() {
	b.angle = 180 - b.angle - width + 20*rand.Float64() // nolint:gosec
}

func (b *Local) increaseSpeed() {
	if b.bounces%2 != 0 {
		return
	}

	switch b.level {
	case level.Easy:
		b.speed += 0.5
	case level.Medium:
		b.speed++
	case level.Hard:
		b.speed += 2
	}

	if b.speed > maxSpeed {
		b.speed = maxSpeed
		return
	}
}
