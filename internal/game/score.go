package game

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font" // Your custom font package
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// score represents the player's score.
type score struct {
	value    int8
	position geometry.Vector
	textFace *text.GoTextFace
}

// newScore creates a new score instance.
func newScore(fontLoader *font.Font, posX float64) *score {
	textFace, err := fontLoader.Face("score", 60)
	if err != nil {
		panic(err)
	}
	return &score{
		value:    0,
		position: geometry.Vector{X: posX, Y: 80},
		textFace: textFace,
	}
}

// draw renders the score on the screen.
func (s *score) draw(screen *ebiten.Image) {
	uiText := ui.Text{
		Value:    strconv.Itoa(int(s.value)),
		FontFace: s.textFace,
		Position: s.position,
		Color:    color.RGBA{200, 200, 200, 255},
	}
	uiText.Draw(screen)
}
