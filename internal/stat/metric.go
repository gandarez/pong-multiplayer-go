package metric

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// nolint: gochecknoglobals
var defaultColor = color.RGBA{0, 0, 0, 255}

// Metric represents the game metric.
type Metric struct {
	screenWidth int
	textFace    *text.GoTextFace
}

// New creates a new metric.
func New(font *font.Font, screenWidth int) (*Metric, error) {
	textFace, err := font.Face("stat", 9)
	if err != nil {
		return nil, fmt.Errorf("failed to create text face: %w", err)
	}

	return &Metric{
		screenWidth: screenWidth,
		textFace:    textFace,
	}, nil
}

// Draw draws the metric on the screen.
// bounces is the number of bounces of the ball.
// angle is the angle of the ball.
// lvl is the current level.
func (m *Metric) Draw(screen *ebiten.Image, bounces int, angle float64, lvl level.Level) {
	// draw current FPS, bounces, current level, or ping if present
	fpsText := fmt.Sprintf(
		"FPS: %.f | Ball (bounces: %d, angle: %.f ) | Level: %s",
		ebiten.ActualFPS(),
		bounces,
		angle,
		lvl.String(),
	)

	uiText := ui.Text{
		Value:    fpsText,
		FontFace: m.textFace,
		Position: geometry.Vector{
			X: 5,
			Y: 0,
		},
		Color: defaultColor,
	}

	uiText.Draw(screen)
}

// DrawNetworkInfo draws the network information on the screen.
func (m *Metric) DrawNetworkInfo(screen *ebiten.Image, pingCurrentPlayer, pingOpponent int64) {
	pingText := fmt.Sprintf(
		"me: %4dms | opponent: %4dms",
		pingCurrentPlayer, pingOpponent,
	)

	width, _ := text.Measure(pingText, m.textFace, 1)

	uiText := ui.Text{
		Value:    pingText,
		FontFace: m.textFace,
		Position: geometry.Vector{
			X: float64(m.screenWidth) - width - 5,
			Y: 0,
		},
		Color: defaultColor,
	}

	uiText.Draw(screen)
}
