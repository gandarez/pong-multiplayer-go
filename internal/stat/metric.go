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

type Metric struct {
	textFace *text.GoTextFace
}

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

func (m *Metric) Draw(screen *ebiten.Image, bounces int, lvl level.Level) {
	// draw current FPS, bounces, current level
	fpsText := fmt.Sprintf(
		"FPS: %.f | Bounces: %d | Level: %s",
		ebiten.ActualFPS(),
		bounces,
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
