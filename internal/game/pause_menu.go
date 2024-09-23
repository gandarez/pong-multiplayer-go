package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font" // Your custom font package
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// PauseMenu represents the pause menu.
type PauseMenu struct {
	font          *font.Font
	screenWidth   float64
	options       []string
	selectedIndex int
	ShouldExit    bool
	ShouldResume  bool
	isOpen        bool
}

// NewPauseMenu creates a new PauseMenu.
func NewPauseMenu(fontLoader *font.Font, screenWidth float64) *PauseMenu {
	return &PauseMenu{
		font:        fontLoader,
		screenWidth: screenWidth,
		options:     []string{"Resume", "Exit"},
	}
}

// Open opens the pause menu.
func (pm *PauseMenu) Open() {
	pm.isOpen = true
	pm.ShouldExit = false
	pm.ShouldResume = false
	pm.selectedIndex = 0
}

// Update updates the pause menu.
func (pm *PauseMenu) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if pm.selectedIndex > 0 {
			pm.selectedIndex--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if pm.selectedIndex < len(pm.options)-1 {
			pm.selectedIndex++
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch pm.selectedIndex {
		case 0:
			pm.ShouldResume = true
			pm.isOpen = false
		case 1:
			pm.ShouldExit = true
		}
	}
}

// Draw draws the pause menu.
func (pm *PauseMenu) Draw(screen *ebiten.Image) {
	// Reduce alpha of the background
	overlay := ebiten.NewImage(int(ScreenWidth), int(ScreenHeight))
	overlay.Fill(ui.TransparentBlack)
	screen.DrawImage(overlay, nil)

	// Draw menu options
	textFace, err := pm.font.Face("ui", 20)
	if err != nil {
		panic(err)
	}
	y := ScreenHeight/2 - float64(len(pm.options)*30)/2

	for i, option := range pm.options {
		color := ui.DefaultColor
		if i == pm.selectedIndex {
			color = ui.HighlightColor
		}
		textWidth, _ := text.Measure(option, textFace, 1)
		uiText := ui.Text{
			Value:    option,
			FontFace: textFace,
			Position: geometry.Vector{
				X: (pm.screenWidth - textWidth) / 2,
				Y: y + float64(i*30),
			},
			Color: color,
		}
		uiText.Draw(screen)
	}
}
