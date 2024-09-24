package game

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// score represents the player's score.
type score struct {
	value    int8
	position geometry.Vector
	textFace *text.GoTextFace
}

func new(textFace *text.GoTextFace, posX float64) *score {
	return &score{
		value:    0,
		position: geometry.Vector{X: posX, Y: 30},
		textFace: textFace,
	}
}

// newScore1 creates a new score instance.
func newScore1(font *font.Font) *score {
	textFace, err := font.Face("score", 60)
	if err != nil {
		panic(err)
	}

	scoreWidth, _ := text.Measure("0", textFace, 1)

	return new(textFace, ScreenWidth/2-50-scoreWidth)
}

// newScore2 creates a new score instance.
func newScore2(font *font.Font) *score {
	textFace, err := font.Face("score", 60)
	if err != nil {
		panic(err)
	}

	return new(textFace, ScreenWidth/2+70)
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
