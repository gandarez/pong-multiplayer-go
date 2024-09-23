package menu

import "github.com/hajimehoshi/ebiten/v2"

// Screen is the interface that wraps the DrawImage method.
type Screen interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions)
}

// State is the interface that represents a state.
type State interface {
	Update()
	Draw(screen Screen)
	String() string
}
