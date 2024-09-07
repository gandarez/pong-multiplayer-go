package player

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	bouncerHeight = 50
	bouncerWidth  = 10
	movementSpeed = 4
)

type (
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

	// , &Player{
	// 	name:          player2Name,
	// 	bouncerHeight: bouncerHeight,
	// 	bouncerWidth:  bouncerWidth,
	// 	position: geometry.Vector{
	// 		X: screenWidth - 25,
	// 		Y: (screenHeight - bouncerHeight) / 2,
	// 	},
	// 	screenHeight:     screenHeight,
	// 	screenWidth:      screenWidth,
	// 	fieldBorderWidth: fieldBorderWidth,
	// }
}

func (p *Player) Update(up, down ebiten.Key) {
	if ebiten.IsKeyPressed(up) {
		if p.position.Y <= p.fieldBorderWidth {
			return
		}

		p.position.Y -= movementSpeed
		return
	}

	if ebiten.IsKeyPressed(down) {
		if p.position.Y > p.screenHeight-p.bouncerHeight-p.fieldBorderWidth {
			return
		}

		p.position.Y += movementSpeed
		return
	}
}

func (p *Player) Bounds() geometry.Rect {
	return geometry.Rect{
		X:      p.position.X,
		Y:      p.position.Y,
		Width:  p.bouncerWidth,
		Height: p.bouncerHeight,
	}
}

func (p *Player) BouncerWidth() float64 {
	return p.bouncerWidth
}

func (p *Player) BouncerHeight() float64 {
	return p.bouncerHeight
}

func (p *Player) Position() geometry.Vector {
	return p.position
}
