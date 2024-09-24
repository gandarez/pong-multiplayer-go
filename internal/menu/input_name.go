package menu

import (
	"log/slog"
	"math"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const maxNameLength = 10

var validNameRegexp = regexp.MustCompile(`^[a-zA-Z-\.]+$`)

// InputNameState is the state where the player can input their name.
type InputNameState struct {
	cursorVisible bool
	cursorTicker  *time.Ticker
	menu          *Menu
}

var _ State = (*InputNameState)(nil)

// NewInputNameState creates a new InputNameState.
func NewInputNameState(menu *Menu) *InputNameState {
	state := &InputNameState{
		menu: menu,
	}

	state.cursorTicker = time.NewTicker(500 * time.Millisecond)
	go func() {
		for range state.cursorTicker.C {
			state.cursorVisible = !state.cursorVisible
		}
	}()

	return state
}

// Update updates the state.
func (s *InputNameState) Update() {
	for _, char := range ebiten.AppendInputChars(nil) {
		if !validNameRegexp.MatchString(string(char)) {
			continue
		}

		//  limit the name length
		if len(s.menu.playerName) == maxNameLength {
			continue
		}

		// do not allow starting with dot or dash
		if len(s.menu.playerName) == 0 && (char == '.' || char == '-') {
			continue
		}

		// do not allow consecutive dots or dashes
		if len(s.menu.playerName) > 0 {
			lastChar, _ := utf8.DecodeLastRuneInString(s.menu.playerName)
			if (lastChar == '.' || lastChar == '-') && (char == '.' || char == '-') {
				continue
			}
		}

		s.menu.playerName += string(char)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(s.menu.playerName) > 0 {
		_, size := utf8.DecodeLastRuneInString(s.menu.playerName)
		s.menu.playerName = s.menu.playerName[:len(s.menu.playerName)-size]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(s.menu.playerName) > 0 {
		// trim any dot or dash at the end
		lastChar, _ := utf8.DecodeLastRuneInString(s.menu.playerName)
		if lastChar == '.' || lastChar == '-' {
			s.menu.playerName = s.menu.playerName[:len(s.menu.playerName)-1]
		}

		s.menu.gameMode = Multiplayer
		s.menu.level = level.Medium
		s.menu.readyToPlay = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.playerName = ""
		s.menu.ChangeState(NewMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *InputNameState) Draw(screen *ebiten.Image) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))
		return
	}

	prompt := "Enter your name:"
	width, _ := text.Measure(prompt, textFace, 1)
	y := 250.0
	uiText := ui.Text{
		Value:    prompt,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (float64(s.menu.screenWidth) - width) / 2,
			Y: y,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)

	name := s.menu.playerName

	widthCursor, _ := text.Measure(name+"_", textFace, 1)
	widthWithoutCursor, _ := text.Measure(name+"", textFace, 1)

	// get max width to not make text jump
	widthName := math.Max(widthCursor, widthWithoutCursor)

	// only show cursor if the name is not at max length
	if s.cursorVisible && len(name) < maxNameLength {
		name += "_"
	}

	uiText = ui.Text{
		Value:    name,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (float64(s.menu.screenWidth) - widthName) / 2,
			Y: y + 30,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)
}

// String returns the state name.
func (*InputNameState) String() string {
	return "InputNameState"
}
