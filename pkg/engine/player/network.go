package player

import "github.com/gandarez/pong-multiplayer-go/pkg/geometry"

// Network represents a player that is controlled by the network.
// It is used to represent the player on the other side of the network.
type Network struct {
	*player
}

// NewNetwork creates a new player to play in a network game.
func NewNetwork(name string, side geometry.Side, screenWidth, screenHeight float64) *Network {
	return &Network{
		player: &player{
			name:          name,
			side:          side,
			bouncerHeight: bouncerHeight,
			bouncerWidth:  bouncerWidth,
			position: geometry.Vector{
				X: calculatePositionX(side, screenWidth),
				Y: (screenHeight - bouncerHeight) / 2,
			},
		},
	}
}

// BouncerHeight returns the height of the bouncer.
func (p *Network) BouncerHeight() float64 {
	return p.bouncerHeight
}

// BouncerWidth returns the width of the bouncer.
func (p *Network) BouncerWidth() float64 {
	return p.bouncerWidth
}

// Bounds returns the bounds of the player.
func (p *Network) Bounds() geometry.Rect {
	return p.player.Bounds()
}

// Name returns the name of the player.
func (p *Network) Name() string {
	return p.name
}

// Position returns the position of the player.
func (p *Network) Position() geometry.Vector {
	return p.position
}

// Side returns the side of the player.
func (p *Network) Side() geometry.Side {
	return p.side
}

// Reset will panic because it is not implemented.
func (*Network) Reset() {
	panic("not implemented")
}

// SetPosition sets the Y position of the player.
func (p *Network) SetPosition(y float64) {
	p.position.Y = y
}

// Update will panic because it is not implemented.
func (*Network) Update(_ Input) {
	panic("not implemented")
}
