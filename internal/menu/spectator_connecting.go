package menu

import (
	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// spectatorConnectingState is the state that represents the connection to a game as a spectator.
type spectatorConnectingState struct {
	menu *Menu
}

// newSpectatorConnectingState creates a new spectatorConnectingState.
func newSpectatorConnectingState(menu *Menu) *spectatorConnectingState {
	return &spectatorConnectingState{
		menu: menu,
	}
}

// Update updates the state.
func (s *spectatorConnectingState) Update() {
	// Set the game mode to Spectator
	s.menu.gameMode = Spectator
	s.menu.readyToPlay = true
}

// Draw draws the state.
func (s *spectatorConnectingState) Draw(screen *ebiten.Image) {
	s.drawConnectingMessage(screen, s.menu.font, float64(s.menu.screenWidth))
}

// String returns the string representation of the state.
func (*spectatorConnectingState) String() string {
	return "spectatorConnectingState"
}

func (*spectatorConnectingState) drawConnectingMessage(screen *ebiten.Image, font *font.Font, screenWidth float64) {
	textFace, err := font.Face("ui", 20)
	if err != nil {
		panic(err)
	}

	message := "Connecting..."
	width, _ := text.Measure(message, textFace, 1)

	uiText := ui.Text{
		Value:    message,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (screenWidth - width) / 2,
			Y: 250.0,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)
}
