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

	player struct {
		name          string
		side          geometry.Side
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
		Side() geometry.Side
		Position() geometry.Vector
		Reset()
		SetPosition(y float64)
		Update(input Input)
	}
)

// calculatePositionX calculates the initial X position of the player.
func calculatePositionX(side geometry.Side, screenWidth float64) float64 {
	x := 15.0
	if side == geometry.Right {
		x = screenWidth - 25
	}

	return x
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
