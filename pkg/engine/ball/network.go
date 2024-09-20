package ball

import "github.com/gandarez/pong-multiplayer-go/pkg/geometry"

// Network represents the ball in a network game.
type Network struct {
	*ball
}

func newNetwork() *Network {
	return &Network{
		ball: &ball{
			width: width,
		},
	}
}

// Angle will panic because it is not implemented.
func (*Network) Angle() float64 {
	panic("not implemented")
}

// Bounces will panic because it is not implemented.
func (*Network) Bounces() int {
	panic("not implemented")
}

// Bounds returns the bounds of the ball.
func (b *Network) Bounds() geometry.Rect {
	return b.ball.Bounds()
}

// CheckGoal will panic because it is not implemented.
func (*Network) CheckGoal() (bool, geometry.Side) {
	panic("not implemented")
}

// Position returns the position of the ball.
func (b *Network) Position() geometry.Vector {
	return b.position
}

// Reset will panic because it is not implemented.
func (*Network) Reset(_ geometry.Side) Ball {
	panic("not implemented")
}

// SetPosition sets the position of the ball.
func (b *Network) SetPosition(pos geometry.Vector) {
	b.position = pos
}

// Update will panic because it is not implemented.
func (*Network) Update(_, _ geometry.Rect) {
	panic("not implemented")
}

// Width returns the width of the ball.
func (b *Network) Width() float64 {
	return b.width
}
