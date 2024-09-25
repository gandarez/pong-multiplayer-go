package ball

import (
	"math"
	"math/rand/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// Local represents the ball in a local game.
type Local struct {
	level        level.Level
	nextSide     geometry.Side
	screenHeight float64
	screenWidth  float64
	speed        float64
	*ball
}

// NewLocal creates a new ball to play locally.
// screenWidth and screenHeight are the dimensions of the screen.
// lvl is the level of the game.
func NewLocal(screenWidth, screenHeight float64, lvl level.Level) *Local {
	var nextSide geometry.Side
	if rand.IntN(2) == 0 { // nolint:gosec
		nextSide = geometry.Left
	} else {
		nextSide = geometry.Right
	}

	return &Local{
		level:        lvl,
		nextSide:     nextSide,
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
func (b *Local) Reset() Ball {
	if b.nextSide == geometry.Left {
		b.nextSide = geometry.Right
	} else {
		b.nextSide = geometry.Left
	}

	return NewLocal(b.screenWidth, b.screenHeight, b.level)
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
		return -45 + float64(rand.IntN(91)) // nolint:gosec
	}

	return 135 + float64(rand.IntN(91)) // nolint:gosec
}

// bounce checks if the ball is bouncing on the walls or the players and changes the angle of the ball.
func (b *Local) bounce(p1Bounds, p2Bounds geometry.Rect) {
	b.checkWallBounce()
	b.checkPaddleBounce(p1Bounds, p2Bounds)
}

// checkWallBounce checks if the ball bounces off the top or bottom walls.
func (b *Local) checkWallBounce() {
	if b.position.Y <= width {
		b.position.Y = width
		b.bounceOffWall()
	}

	if b.position.Y >= b.screenHeight-b.width-width {
		b.position.Y = b.screenHeight - b.width - width
		b.bounceOffWall()
	}
}

// bounceOffWall changes the ball's angle when it hits a wall and slightly adjusts its angle randomly.
func (b *Local) bounceOffWall() {
	b.bounces++
	b.angle *= -1
	// nolint:gosec
	// slight random adjustment to avoid flat bounces
	b.angle += 5 * (rand.Float64() - 0.5)
	b.increaseSpeed()
}

// checkPaddleBounce checks if the ball is hitting one of the paddles and bounces off.
func (b *Local) checkPaddleBounce(p1Bounds, p2Bounds geometry.Rect) {
	if b.ball.Bounds().Intersects(p1Bounds) {
		b.position.X = p1Bounds.X + p1Bounds.Width + width
		b.bounceOffPaddle()
	}

	if b.ball.Bounds().Intersects(p2Bounds) {
		b.position.X = p2Bounds.X - b.width
		b.bounceOffPaddle()
	}
}

// bounceOffPaddle changes the ball's angle when it hits a paddle and slightly randomizes the angle.
func (b *Local) bounceOffPaddle() {
	b.bounces++
	b.randomBounce()
	b.increaseSpeed()
}

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
