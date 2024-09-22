package menu

import (
	"log/slog"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	// menu strings.
	titleStr        = "PONGO"
	backStr         = "Back"
	localModeStr    = "Local Mode"
	onePlayerStr    = "One Player"
	twoPlayersStr   = "Two Players"
	multiplayerStr  = "Multiplayer"
	instructionsStr = "Instructions"
)

// state is the state of the menu.
type state int

const (
	mainMenu state = iota
	localMode
	levelSelection
	instructions
)

// GameMode is the game mode.
type GameMode int

const (
	// Undefined represents an undefined/unset game mode.
	Undefined GameMode = iota
	// OnePlayer represents a single player game mode.
	OnePlayer
	// TwoPlayers represents a two players game mode.
	TwoPlayers
	// Multiplayer represents a multiplayer game mode.
	Multiplayer
)

// Menu represents the game menu.
type Menu struct {
	font           *font.Font
	gameMode       GameMode
	level          level.Level
	state          state
	selectedOption int8
	screenWidth    int
	readyToPlay    bool
}

// New creates a new game menu.
func New(font *font.Font, screenWidth int) *Menu {
	return &Menu{
		font:        font,
		gameMode:    Undefined,
		state:       mainMenu,
		screenWidth: screenWidth,
	}
}

// Update updates the menu state.
// nolint:gocyclo
func (m *Menu) Update() {
	// key down
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		switch m.state {
		case mainMenu:
			switch m.selectedOption {
			case 0: // local mode
				m.selectedOption = 1
			case 1: // multiplayer mode
				m.selectedOption = 2 // instructions
			}
		case localMode:
			switch m.selectedOption {
			case 0: // one player
				m.selectedOption = 1
			case 1: // two players
				m.selectedOption = 2 // back
			}
		case levelSelection:
			switch m.selectedOption {
			case 0: // easy
				m.selectedOption = 1
			case 1: // medium
				m.selectedOption = 2
			case 2: // hard
				m.selectedOption = 3 // back
			}
		}

		return
	}

	// key up
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		switch m.state {
		case mainMenu:
			switch m.selectedOption {
			case 2: // instructions
				m.selectedOption = 1
			case 1: // multiplayer mode
				m.selectedOption = 0 // local mode
			}
		case localMode:
			switch m.selectedOption {
			case 2: // back
				m.selectedOption = 1
			case 1: // two players
				m.selectedOption = 0 // one player
			}
		case levelSelection:
			switch m.selectedOption {
			case 3: // back
				m.selectedOption = 2
			case 2: // hard
				m.selectedOption = 1
			case 1: // medium
				m.selectedOption = 0 // easy
			}
		}

		return
	}

	// key enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.state {
		case mainMenu:
			switch m.selectedOption {
			case 0:
				// local mode
				m.state = localMode
				m.selectedOption = 0
			case 1:
				// multiplayer mode
				m.gameMode = Multiplayer
				m.state = levelSelection
				m.selectedOption = 0
			case 2:
				// instructions
				m.state = instructions
				m.selectedOption = 0
			}
		case localMode:
			switch m.selectedOption {
			case 0:
				// single player
				m.gameMode = OnePlayer
				m.state = levelSelection
				m.selectedOption = 0
			case 1:
				// two players
				m.gameMode = TwoPlayers
				m.state = levelSelection
				m.selectedOption = 0
			case 2:
				// back
				m.state = mainMenu
				m.gameMode = Undefined
				m.readyToPlay = false
				m.selectedOption = 0
			}
		case levelSelection:
			switch m.selectedOption {
			case 0:
				// easy
				m.level = level.Easy
				m.readyToPlay = true
			case 1:
				// medium
				m.level = level.Medium
				m.readyToPlay = true
			case 2:
				// hard
				m.level = level.Hard
				m.readyToPlay = true
			case 3:
				// back
				m.state = mainMenu
				m.gameMode = Undefined
				m.readyToPlay = false
				m.selectedOption = 0
			}
		}
	}

	// key esc
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch m.state {
		case mainMenu:
			// exit the game
			m.gameMode = Undefined
			m.readyToPlay = true // just to force exit menu
		case localMode, levelSelection, instructions:
			m.state = mainMenu
			m.selectedOption = 0
			m.gameMode = Undefined
			m.readyToPlay = false
		}
	}
}

// Draw draws the menu on the screen.
func (m *Menu) Draw(screen *ebiten.Image) {
	switch m.state {
	case mainMenu:
		m.drawOption(screen, localModeStr, multiplayerStr, instructionsStr)
	case localMode:
		m.drawOption(screen, onePlayerStr, twoPlayersStr, backStr)
	case levelSelection:
		m.drawOption(screen, level.Easy.String(), level.Medium.String(), level.Hard.String(), backStr)
	case instructions:
		m.drawInstructions(screen)
	}
}

func (m *Menu) drawOption(screen *ebiten.Image, values ...string) {
	textFace, err := m.font.Face("ui", 20)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))

		return
	}

	var maxWidth float64

	y := 250.0

	for _, val := range values {
		width, _ := text.Measure(val, textFace, 1)

		uiText := ui.Text{
			Value:    val,
			FontFace: textFace,
			Position: geometry.Vector{
				X: (float64(m.screenWidth) - width) / 2,
				Y: y,
			},
			Color: ui.DefaultColor,
		}

		uiText.Draw(screen)

		maxWidth = math.Max(maxWidth, width)

		y += 50
	}

	// draw selected option
	y = 255. + 50*float64(m.selectedOption)

	vector.DrawFilledRect(
		screen,
		float32(float64(m.screenWidth)-maxWidth)/2-30, float32(y),
		15, 15, ui.DefaultColor, true,
	)
}

func (m *Menu) drawInstructions(screen *ebiten.Image) {
	textFace, err := m.font.Face("ui", 12)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))

		return
	}

	var maxWidth float64

	y := 200.0

	val := strings.ReplaceAll(instructionsDetailedStr, "\r\n", "\n")
	splitted := strings.Split(val, "\n")

	var previousStr string

	for _, str := range splitted {
		width, height := text.Measure(str, textFace, 1)
		if previousStr == "" && height == 0 {
			height = 10
		}

		uiText := ui.Text{
			Value:    str,
			FontFace: textFace,
			Position: geometry.Vector{
				X: (float64(m.screenWidth) - width) / 2,
				Y: y,
			},
			Color: ui.DefaultColor,
		}

		uiText.Draw(screen)

		y += float64(height)

		maxWidth = math.Max(maxWidth, width)
	}

	// draw selected option
	y = 255. + 50*float64(m.selectedOption)

	vector.DrawFilledRect(
		screen,
		float32(float64(m.screenWidth)-maxWidth)/2-30, float32(y),
		15, 15, ui.DefaultColor, true,
	)
}

// IsReadyToPlay returns true if the game is ready to play.
func (m *Menu) IsReadyToPlay() bool {
	return m.readyToPlay
}

// GameMode returns the selected game mode.
func (m *Menu) GameMode() GameMode {
	return m.gameMode
}

// Level returns the selected level.
func (m *Menu) Level() level.Level {
	return m.level
}
