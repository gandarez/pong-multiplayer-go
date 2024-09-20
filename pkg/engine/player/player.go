package player

import (
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	bouncerHeight = 50
	bouncerWidth  = 10
	movementSpeed = 4
)

type (
	// Input represents the input of the player.
	Input struct {
		Up   bool
		Down bool
	}

	// Player is the player of the game.
	Player struct {
		name             string
		bouncerHeight    float64
		bouncerWidth     float64
		position         geometry.Vector
		screenWidth      float64
		screenHeight     float64
		fieldBorderWidth float64
	}
)

// TODO: find a better way to handle multiplayer players.
func NewMultiplayer(name string, side geometry.Side, screenWidth, screenHeight float64) *Player {
	x := 15.0
	if side == geometry.Right {
		x = screenWidth - 25
	}

	return &Player{
		bouncerHeight: bouncerHeight,
		bouncerWidth:  bouncerWidth,
		name:          name,
		position: geometry.Vector{
			X: x,
			Y: (screenHeight - bouncerHeight) / 2,
		},
	}
}

// New creates a new player.
func New(name string, side geometry.Side, screenWidth, screenHeight, fieldBorderWidth float64) *Player {
	x := 15.0
	if side == geometry.Right {
		x = screenWidth - 25
	}

	return &Player{
		name:          name,
		bouncerHeight: bouncerHeight,
		bouncerWidth:  bouncerWidth,
		position: geometry.Vector{
			X: x,
			Y: (screenHeight - bouncerHeight) / 2,
		},
		screenHeight:     screenHeight,
		screenWidth:      screenWidth,
		fieldBorderWidth: fieldBorderWidth,
	}
}

// Reset resets the player to its initial position.
func (p *Player) Reset() {
	p.position.Y = (p.screenHeight - p.bouncerHeight) / 2
}

// Update updates the player position based on the keys pressed.
func (p *Player) Update(input Input) {
	switch {
	case input.Up:
		p.position.Y -= movementSpeed
	case input.Down:
		p.position.Y += movementSpeed
	}

	p.keepInBounds()
}

func (p *Player) keepInBounds() {
	if p.position.Y < p.fieldBorderWidth {
		p.position.Y = p.fieldBorderWidth
	}

	if p.position.Y > p.screenHeight-p.bouncerHeight-p.fieldBorderWidth {
		p.position.Y = p.screenHeight - p.bouncerHeight - p.fieldBorderWidth
	}
}

// Bounds returns the bounds of the player.
func (p *Player) Bounds() geometry.Rect {
	return geometry.Rect{
		X:      p.position.X,
		Y:      p.position.Y,
		Width:  p.bouncerWidth,
		Height: p.bouncerHeight,
	}
}

// BouncerWidth returns the width of the player bouncer.
func (p *Player) BouncerWidth() float64 {
	return p.bouncerWidth
}

// BouncerHeight returns the height of the player bouncer.
func (p *Player) BouncerHeight() float64 {
	return p.bouncerHeight
}

// Position returns the position of the player.
func (p *Player) Position() geometry.Vector {
	return p.position
}

// SetPosition sets the Y position of the player.
func (p *Player) SetPosition(y float64) {
	p.position.Y = y
}

// Name returns the name of the player.
func (p *Player) Name() string {
	return p.name
}
