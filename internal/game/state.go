package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
)

// state represents a state of the game.
type state interface {
	update() error
	draw(screen *ebiten.Image)
	getBall() ball.Ball
	canPause() bool
}
