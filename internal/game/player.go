package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	engineplayer "github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
)

// drawPlayer draws a player on the screen.
func drawPlayer(p engineplayer.Player, screen *ebiten.Image) {
	for x := 0.0; x < p.BouncerWidth(); x++ {
		for y := 0.0; y < p.BouncerHeight(); y++ {
			screen.Set(int(p.Position().X+x), int(p.Position().Y+y), color.RGBA{200, 200, 200, 255})
		}
	}
}
