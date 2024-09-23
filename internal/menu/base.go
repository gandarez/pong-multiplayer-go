package menu

import (
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

type BaseState struct {
	selectedOption int
	options        []string
	menu           *Menu
}

func (s *BaseState) navigateOptions(maxOptions int) {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if s.selectedOption < maxOptions-1 {
			s.selectedOption++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if s.selectedOption > 0 {
			s.selectedOption--
		}
	}
}

func (s *BaseState) drawOptions(screen Screen) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))
		return
	}

	var maxWidth float64
	y := 250.0

	for _, val := range s.options {
		width, _ := text.Measure(val, textFace, 1)
		uiText := ui.Text{
			Value:    val,
			FontFace: textFace,
			Position: geometry.Vector{
				X: (float64(s.menu.screenWidth) - width) / 2,
				Y: y,
			},
			Color: ui.DefaultColor,
		}
		uiText.Draw(screen.(*ebiten.Image))
		maxWidth = math.Max(maxWidth, width)
		y += 50
	}

	// Draw selected option indicator
	y = 255.0 + 50.0*float64(s.selectedOption)
	vector.DrawFilledRect(
		screen.(*ebiten.Image),
		float32((float64(s.menu.screenWidth)-maxWidth)/2-30), float32(y),
		15, 15, ui.DefaultColor, true,
	)
}
