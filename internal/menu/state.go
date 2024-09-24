package menu

import "github.com/hajimehoshi/ebiten/v2"

// State is the interface that represents a state.
type State interface {
	Draw(screen *ebiten.Image)
	String() string
	Update()
}
