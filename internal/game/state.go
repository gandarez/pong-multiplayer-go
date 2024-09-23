package game

import "github.com/hajimehoshi/ebiten/v2"

// GameState represents a state of the game.
type GameState interface {
	Update() error
	Draw(screen *ebiten.Image)
}
