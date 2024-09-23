package game

import (
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	engineplayer "github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

type player struct {
	namePosition geometry.Vector
	engineplayer.Player
}

func (p *player) draw(screen *ebiten.Image) {
	for x := 0.0; x < p.BouncerWidth(); x++ {
		for y := 0.0; y < p.BouncerHeight(); y++ {
			screen.Set(int(p.Position().X+x), int(p.Position().Y+y), color.RGBA{200, 200, 200, 255})
		}
	}
}

func (p *player) drawName(screen *ebiten.Image, font *font.Font) {
	textface, err := font.Face("ui", 12)
	if err != nil {
		slog.Error("failed to get text face to draw player name", slog.Any("error", err))
		return
	}

	t := ui.Text{
		Value:    p.Name(),
		FontFace: textface,
		Position: p.namePosition,
		Color:    color.RGBA{200, 200, 200, 200},
	}

	t.Draw(screen)
}
