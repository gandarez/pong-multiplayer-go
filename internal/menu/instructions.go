package menu

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const instructionsDetailedStr = `Welcome to PONGO!

The goal of the game is to score points by hitting the ball with your bouncer and 
prevent the opponent from hitting it.

The game ends when one of the players reaches 10 points.

Player 1
- Move up: Q
- Move down: A

Player 2
- Move up: Up arrow
- Move down: Down arrow

Select the game mode, level and press Enter to start the game.

Press Esc to go back to the previous menu.`

// InstructionsState is the state where the player can see the game instructions.
type InstructionsState struct {
	menu *Menu
}

var _ State = (*InstructionsState)(nil)

// NewInstructionsState creates a new InstructionsState.
func NewInstructionsState(menu *Menu) *InstructionsState {
	return &InstructionsState{
		menu: menu,
	}
}

// Update updates the state.
func (s *InstructionsState) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(NewMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *InstructionsState) Draw(screen *ebiten.Image) {
	textFace, err := s.menu.font.Face("ui", 12)
	if err != nil {
		return
	}

	y := 200.0
	val := strings.ReplaceAll(instructionsDetailedStr, "\r\n", "\n")
	splitted := strings.Split(val, "\n")

	for _, str := range splitted {
		width, height := text.Measure(str, textFace, 1)
		if str == "" && height == 0 {
			height = 10
		}

		uiText := ui.Text{
			Value:    str,
			FontFace: textFace,
			Position: geometry.Vector{
				X: (float64(s.menu.screenWidth) - width) / 2,
				Y: y,
			},
			Color: ui.DefaultColor,
		}

		uiText.Draw(screen)

		y += float64(height)
	}
}

// String returns the name of the state.
func (*InstructionsState) String() string {
	return "InstructionsState"
}
