package metric

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// Metric represents the game metric.
type Metric struct {
	textFace *text.GoTextFace
}

// New creates a new metric.
func New(assets *assets.Assets) (*Metric, error) {
	textFaceSource, err := assets.NewTextFaceSource("stat")
	if err != nil {
		return nil, fmt.Errorf("failed to create text face source: %w", err)
	}

	textFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   9,
	}

	return &Metric{
		textFace: textFace,
	}, nil
}

// Draw draws the metric on the screen.
// bounces is the number of bounces of the ball.
// angle is the angle of the ball.
// lvl is the current level.
func (m *Metric) Draw(screen *ebiten.Image, bounces int, angle float64, lvl level.Level) {
	// draw current FPS, bounces, current level
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
		Color: color.RGBA{0, 0, 0, 255},
	}

	uiText.Draw(screen)
}
