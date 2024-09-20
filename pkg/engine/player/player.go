package player

import (
	"errors"

	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	bouncerHeight = 50
	bouncerWidth  = 10
	movementSpeed = 4
)

// Kind represents the kind of player.
type Kind int

const (
	// KindLocal represents a local player.
	KindLocal Kind = iota
	// KindNetwork represents a network player.
	KindNetwork
)

type (
	// Input represents the input of the player.
	Input struct {
		Up   bool
		Down bool
	}

	player struct {
		name          string
		bouncerHeight float64
		bouncerWidth  float64
		position      geometry.Vector
	}

	// Player represents a player.
	Player interface {
		BouncerHeight() float64
		BouncerWidth() float64
		Bounds() geometry.Rect
		Name() string
		Position() geometry.Vector
		Reset()
		SetPosition(y float64)
		Update(input Input)
	}
)

// New creates a new player.
func New(
	kind Kind,
	name string,
	side geometry.Side,
	screenWidth, screenHeight float64,
	fieldBorderWidth *float64,
) (Player, error) {
	x := 15.0
	if side == geometry.Right {
		x = screenWidth - 25
	}

	var player Player

	switch kind {
	case KindLocal:
		if fieldBorderWidth == nil {
			return nil, errors.New("fieldBorderWidth is required for local player")
		}

		player = newLocal(name, x, screenWidth, screenHeight, *fieldBorderWidth)
	case KindNetwork:
		player = newNetwork(name, x, screenHeight)
	}

	return player, nil
}

// Bounds returns the bounds of the player.
func (p *player) Bounds() geometry.Rect {
	return geometry.Rect{
		X:      p.position.X,
		Y:      p.position.Y,
		Width:  p.bouncerWidth,
		Height: p.bouncerHeight,
	}
}
