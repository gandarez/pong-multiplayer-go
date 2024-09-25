package menu

import "github.com/hajimehoshi/ebiten/v2"

// state is the interface that represents a state.
type state interface {
	Draw(screen *ebiten.Image)
	String() string
	Update()
}
