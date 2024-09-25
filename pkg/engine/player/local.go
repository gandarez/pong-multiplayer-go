package player

import "github.com/gandarez/pong-multiplayer-go/pkg/geometry"

// Local represents a player that is controlled by the user.
type Local struct {
	screenWidth      float64
	screenHeight     float64
	fieldBorderWidth float64
	*player
}

// NewLocal creates a new player to play locally.
func NewLocal(name string, side geometry.Side, screenWidth, screenHeight, fieldBorderWidth float64) *Local {
	return &Local{
		screenHeight:     screenHeight,
		screenWidth:      screenWidth,
		fieldBorderWidth: fieldBorderWidth,
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
func (p *Local) BouncerHeight() float64 {
	return p.bouncerHeight
}

// BouncerWidth returns the width of the bouncer.
func (p *Local) BouncerWidth() float64 {
	return p.bouncerWidth
}

// Bounds returns the bounds of the player.
func (p *Local) Bounds() geometry.Rect {
	return p.player.Bounds()
}

// Name returns the name of the player.
func (p *Local) Name() string {
	return p.name
}

// Position returns the position of the player.
func (p *Local) Position() geometry.Vector {
	return p.position
}

// Side returns the side of the player.
func (p *Local) Side() geometry.Side {
	return p.side
}

// Reset resets the position of the player.
func (p *Local) Reset() {
	p.position.Y = (p.screenHeight - p.bouncerHeight) / 2
}

// SetPosition sets the Y position of the player.
func (p *Local) SetPosition(y float64) {
	p.position.Y = y
}

// SetName sets the name of the player.
func (p *Local) SetName(name string) {
	p.name = name
}

// Update updates the position of the player based on the input.
func (p *Local) Update(input Input) {
	switch {
	case input.Up:
		p.position.Y -= movementSpeed
	case input.Down:
		p.position.Y += movementSpeed
	}

	p.keepInBounds()
}

func (p *Local) keepInBounds() {
	if p.position.Y < p.fieldBorderWidth {
		p.position.Y = p.fieldBorderWidth
	}

	if p.position.Y > p.screenHeight-p.bouncerHeight-p.fieldBorderWidth {
		p.position.Y = p.screenHeight - p.bouncerHeight - p.fieldBorderWidth
	}
}
