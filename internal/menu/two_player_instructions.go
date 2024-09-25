package menu

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	instructionsTitle = "Controls"
	lineHeight        = 30
)

// Control defines controls for each player.
type Control struct {
	action string
	key    string
}

// twoPlayersInstructionsState is the state where the controls for two players are displayed.
type twoPlayersInstructionsState struct {
	menu       *Menu
	p1Controls []Control
	p2Controls []Control
}

// newTwoPlayersInstructionsState creates a new twoPlayersInstructionsState.
func newTwoPlayersInstructionsState(menu *Menu) *twoPlayersInstructionsState {
	return &twoPlayersInstructionsState{
		menu: menu,
		p1Controls: []Control{
			{action: "Up", key: "Q"},
			{action: "Down", key: "A"},
		},
		p2Controls: []Control{
			{action: "Up", key: "Up Arrow"},
			{action: "Down", key: "Down Arrow"},
		},
	}
}

// Update handles the logic for the TwoPlayersInstructionsState.
func (s *twoPlayersInstructionsState) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// proceed to level selection after showing the instructions.
		s.menu.ChangeState(newLevelSelectionState(s.menu))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// go back to the LocalModeState.
		s.menu.ChangeState(newLocalModeState(s.menu))
	}
}

// Draw renders the instructions on the screen.
func (s *twoPlayersInstructionsState) Draw(screen *ebiten.Image) {
	// draw the title
	titlePosition, err := s.drawTitle(screen)
	if err != nil {
		slog.Error("failed to draw title", slog.Any("error", err))
	}

	// draw player controls
	err = s.drawControls(screen, titlePosition)
	if err != nil {
		slog.Error("failed to draw controls", slog.Any("error", err))
	}

	// draw instructions at the bottom
	err = s.drawMenuOptions(screen)
	if err != nil {
		slog.Error("failed to draw menu options", slog.Any("error", err))
	}
}

func (s *twoPlayersInstructionsState) drawTitle(screen *ebiten.Image) (geometry.Vector, error) {
	titleFace, err := s.menu.font.Face("ui", 30)
	if err != nil {
		return geometry.Vector{}, err
	}

	titleWidth, _ := text.Measure(instructionsTitle, titleFace, 1)
	titlePosition := geometry.Vector{
		X: (float64(s.menu.screenWidth) - titleWidth) / 2,
		Y: 200,
	}

	uiText := ui.Text{
		Value:    instructionsTitle,
		FontFace: titleFace,
		Position: titlePosition,
		Color:    ui.DefaultColor,
	}
	uiText.Draw(screen)

	return titlePosition, nil
}

func (s *twoPlayersInstructionsState) drawControls(screen *ebiten.Image, titlePosition geometry.Vector) error {
	// define positions for columns
	columnWidth := float64(s.menu.screenWidth) / 2
	leftColumnX := columnWidth / 2
	rightColumnX := columnWidth + leftColumnX

	controlsFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		return err
	}

	// draw player 1 name
	playerNameY := titlePosition.Y + 60
	s.drawPlayerName(screen, controlsFace, "Player 1", leftColumnX, playerNameY)

	// draw player 2 name
	s.drawPlayerName(screen, controlsFace, "Player 2", rightColumnX, playerNameY)

	// draw controls for each player
	startY := playerNameY + 40
	s.drawPlayerControls(screen, controlsFace, s.p1Controls, leftColumnX, startY)
	s.drawPlayerControls(screen, controlsFace, s.p2Controls, rightColumnX, startY)

	return nil
}

func (*twoPlayersInstructionsState) drawPlayerName(
	screen *ebiten.Image,
	controlsFace text.Face,
	playerName string,
	x, y float64,
) {
	playerNameWidth, _ := text.Measure(playerName, controlsFace, 1)
	playerNamePosition := geometry.Vector{
		X: x - playerNameWidth/2,
		Y: y,
	}

	uiText := ui.Text{
		Value:    playerName,
		FontFace: controlsFace,
		Position: playerNamePosition,
		Color:    ui.HighlightColor,
	}
	uiText.Draw(screen)
}

func (*twoPlayersInstructionsState) drawPlayerControls(
	screen *ebiten.Image,
	controlsFace text.Face,
	controls []Control,
	startX, startY float64,
) {
	for i, control := range controls {
		y := startY + float64(i)*lineHeight
		action := control.action
		key := control.key
		controlText := fmt.Sprintf("%s: [%s]", action, key)
		controlWidth, _ := text.Measure(controlText, controlsFace, 1)
		controlPosition := geometry.Vector{
			X: startX - controlWidth/2,
			Y: y,
		}

		uiText := ui.Text{
			Value:    controlText,
			FontFace: controlsFace,
			Position: controlPosition,
			Color:    ui.DefaultColor,
		}
		uiText.Draw(screen)
	}
}

func (s *twoPlayersInstructionsState) drawMenuOptions(screen *ebiten.Image) error {
	instructionsFace, err := s.menu.font.Face("ui", 16)
	if err != nil {
		return err
	}

	// define separate texts
	enterText := "Press Enter to continue"
	escText := "Press Esc to go back"

	// measure the width of each text
	enterTextWidth, _ := text.Measure(enterText, instructionsFace, 1)
	escTextWidth, _ := text.Measure(escText, instructionsFace, 1)

	// calculate text center positions
	centerX := float64(s.menu.screenWidth) / 2
	baseY := float64(s.menu.screenHeight) - 80
	enterTextPosition := geometry.Vector{
		X: centerX - (enterTextWidth / 2),
		Y: baseY,
	}

	escTextPosition := geometry.Vector{
		X: centerX - (escTextWidth / 2),
		Y: baseY + 30,
	}

	// draw "Press Enter" text
	uiText := ui.Text{
		Value:    enterText,
		FontFace: instructionsFace,
		Position: enterTextPosition,
		Color:    ui.DefaultColor,
	}
	uiText.Draw(screen)

	// draw "Press Esc" text
	uiText = ui.Text{
		Value:    escText,
		FontFace: instructionsFace,
		Position: escTextPosition,
		Color:    ui.DefaultColor,
	}
	uiText.Draw(screen)

	return nil
}

// String returns the name of the state.
func (*twoPlayersInstructionsState) String() string {
	return "twoPlayersInstructionsState"
}
