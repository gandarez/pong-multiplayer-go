package menu

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions)
}

type State interface {
	Update()
	Draw(screen Screen)
}
