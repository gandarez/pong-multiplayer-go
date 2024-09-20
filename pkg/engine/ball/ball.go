package ball

import (
	"math"
	"math/rand"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	initialSpeed = 2
	maxSpeed     = 8
	width        = 10
)

// Ball is the ball of the game.
type Ball struct {
	angle        float64
	bounces      int
	level        level.Level
	position     geometry.Vector
	screenHeight float64
	screenWidth  float64
	speed        float64
	width        float64
}

// TODO: find a better way to handle multiplayer players.
func NewMultiplayer() *Ball {
	return &Ball{
		width: width,
	}
}

// New creates a new ball.
// nextSide is the side that the ball will go when the game starts.
// screenWidth and screenHeight are the dimensions of the screen.
func New(nextSide geometry.Side, screenWidth, screenHeight float64, lvl level.Level) *Ball {
	return &Ball{
		angle: calcInitialAngle(nextSide),
		level: lvl,
		position: geometry.Vector{
			X: (screenWidth - width) / 2,
			Y: (screenHeight - width) / 2,
		},
		screenHeight: screenHeight,
		screenWidth:  screenWidth,
		speed:        initialSpeed,
		width:        width,
	}
}

// Reset resets the ball to the center of the screen and changes the side that the ball will go.
func (b *Ball) Reset(nextSide geometry.Side) *Ball {
	return New(nextSide, b.screenWidth, b.screenHeight, b.level)
}

// Bounds returns the bounds of the ball.
func (b *Ball) Bounds() geometry.Rect {
	return geometry.Rect{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  b.width,
		Height: b.width,
	}
}

// Update updates the ball position and checks if it is bouncing on the walls or the players.
func (b *Ball) Update(p1Bounds, p2Bounds geometry.Rect) {
	b.position.X += b.speed * math.Cos(b.angle*math.Pi/180)
	b.position.Y += b.speed * math.Sin(b.angle*math.Pi/180)

	b.bounce(p1Bounds, p2Bounds)
}

// bounce checks if the ball is bouncing on the walls or the players and changes the angle of the ball.
func (b *Ball) bounce(p1Bounds, p2Bounds geometry.Rect) {
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

func (b *Ball) randomBounce() {
	b.angle = 180 - b.angle - width + 20*rand.Float64() // nolint:gosec
}

func (b *Ball) increaseSpeed() {
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

// CheckGoal checks if the ball is out of the field and returns true if it is.
// It also returns the side that ball went out.
func (b *Ball) CheckGoal() (bool, geometry.Side) {
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

// Angle returns the angle of the ball.
func (b *Ball) Angle() float64 {
	return b.angle
}

// Bounces returns the number of bounces that the ball has.
func (b *Ball) Bounces() int {
	return b.bounces
}

// Width returns the width of the ball.
func (b *Ball) Width() float64 {
	return b.width
}

// Position returns the position of the ball.
func (b *Ball) Position() geometry.Vector {
	return b.position
}

// SetPosition sets the position of the ball.
func (b *Ball) SetPosition(pos geometry.Vector) {
	b.position = pos
}

func calcInitialAngle(nextSide geometry.Side) float64 {
	if nextSide == geometry.Left {
		return -45 + float64(rand.Intn(91)) // nolint:gosec
	}

	return 135 + float64(rand.Intn(91)) // nolint:gosec
}
