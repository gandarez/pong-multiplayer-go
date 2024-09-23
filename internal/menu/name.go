package menu

import (
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

type InputNameState struct {
	menu          *Menu
	cursorVisible bool
	cursorTicker  *time.Ticker
}

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

func (s *InputNameState) Update() {
	for _, char := range ebiten.InputChars() {
		if char == '\n' || char == '\r' {
			continue
		}
		s.menu.playerName += string(char)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(s.menu.playerName) > 0 {
		_, size := utf8.DecodeLastRuneInString(s.menu.playerName)
		s.menu.playerName = s.menu.playerName[:len(s.menu.playerName)-size]
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(s.menu.playerName) > 0 {
		s.menu.gameMode = Multiplayer
		s.menu.level = level.Medium
		s.menu.readyToPlay = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(NewMainMenuState(s.menu))
	}
}

func (s *InputNameState) Draw(screen Screen) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
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
	uiText.Draw(screen.(*ebiten.Image))

	name := s.menu.playerName
	if s.cursorVisible {
		name += "_"
	}
	widthName, _ := text.Measure(name, textFace, 1)
	uiText = ui.Text{
		Value:    name,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (float64(s.menu.screenWidth) - widthName) / 2,
			Y: y + 30,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen.(*ebiten.Image))
}
