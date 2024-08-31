package engine

import (
	"image/color"

	"github.com/gandarez/pong-multiplayer-go/internal/geometry"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerBouncerHeight = 50
	playerBouncerWidth  = 10
	playerMovementSpeed = 4

	player1Name = "Player 1"
	player2Name = "Player 2"
)

type player struct {
	name          string
	bouncerHeight float64
	bouncerWidth  float64
	position      geometry.Vector
}

// newPlayers creates player 1 and player 2.
func newPlayers() (*player, *player) {
	return &player{
			name:          player1Name,
			bouncerHeight: playerBouncerHeight,
			bouncerWidth:  playerBouncerWidth,
			position: geometry.Vector{
				X: 15,
				Y: (ScreenHeight - playerBouncerHeight) / 2,
			},
		}, &player{
			name:          player2Name,
			bouncerHeight: playerBouncerHeight,
			bouncerWidth:  playerBouncerWidth,
			position: geometry.Vector{
				X: ScreenWidth - 25,
				Y: (ScreenHeight - playerBouncerHeight) / 2,
			},
		}
}

func (p *player) draw(screen *ebiten.Image) {
	for x := 0.0; x < p.bouncerWidth; x++ {
		for y := 0.0; y < p.bouncerHeight; y++ {
			screen.Set(int(p.position.X+x), int(p.position.Y+y), color.RGBA{200, 200, 200, 255})
		}
	}
}

func (p *player) update(up, down ebiten.Key) {
	if ebiten.IsKeyPressed(up) {
		if p.position.Y <= fieldWidth {
			return
		}
		p.position.Y -= playerMovementSpeed

		return
	}

	if ebiten.IsKeyPressed(down) {
		if p.position.Y > ScreenHeight-p.bouncerHeight-fieldWidth {
			return
		}
		p.position.Y += playerMovementSpeed

		return
	}
}

func (p *player) bounds() geometry.Rect {
	return geometry.Rect{
		X:      p.position.X,
		Y:      p.position.Y,
		Width:  p.bouncerWidth,
		Height: p.bouncerHeight,
	}
}
