package menu

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	// menu strings
	backStr                    = "Back"
	localModeStr               = "Local Mode"
	multiplayerStr             = "Multiplayer"
	multiplayerFindOpponentStr = "Finding opponent.."
	onePlayerStr               = "One Player"
	titleStr                   = "PONGO"
	twoPlayersStr              = "Two Players"

	// levels
	levelEasyStr   = "Easy"
	levelMediumStr = "Medium"
	levelHardStr   = "Hard"
)

// state is the state of the menu.
type state int

const (
	mainMenu state = iota
	localMode
	multiplayerMode
	levelSelection
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
	mainTitleTextFace *text.GoTextFace
	gameMode          GameMode
	level             level.Level
	state             state
	selectedOption    int8
	screenWidth       float64
	readyToPlay       bool
}

// New creates a new game menu.
func New(assets *assets.Assets, screenWidth int) (*Menu, error) {
	mainTitleTextFaceSource, err := assets.NewTextFaceSource("ui")
	if err != nil {
		return nil, fmt.Errorf("failed to create main title text face source: %w", err)
	}

	mainTitleTextFace := &text.GoTextFace{
		Source: mainTitleTextFaceSource,
		Size:   80,
	}

	return &Menu{
		mainTitleTextFace: mainTitleTextFace,
		gameMode:          Undefined,
		state:             mainMenu,
		screenWidth:       float64(screenWidth),
	}, nil
}

// Update updates the menu state.
func (m *Menu) Update() {
	// key down
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		switch m.state {
		case mainMenu:
			if m.selectedOption == 0 {
				m.selectedOption++
			}
		case localMode:
			switch m.selectedOption {
			case 0:
				m.selectedOption = 1
			case 1:
				m.selectedOption = 2
			}
		case levelSelection:
			switch m.selectedOption {
			case 0:
				m.selectedOption = 1
			case 1:
				m.selectedOption = 2
			case 2:
				m.selectedOption = 3
			}
		}

		return
	}

	// key up
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		switch m.state {
		case mainMenu:
			if m.selectedOption == 1 {
				m.selectedOption--
			}
		case localMode:
			switch m.selectedOption {
			case 2:
				m.selectedOption = 1
			case 1:
				m.selectedOption = 0
			}
		case levelSelection:
			switch m.selectedOption {
			case 3:
				m.selectedOption = 2
			case 2:
				m.selectedOption = 1
			case 1:
				m.selectedOption = 0
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
				m.state = localMode
				m.selectedOption = 0
			case 1:
				m.state = multiplayerMode
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
				m.state = localMode
				m.gameMode = Undefined
				m.readyToPlay = false
				m.selectedOption = 0
			}
		}

		return
	}
}

// Draw draws the menu on the screen.
func (m *Menu) Draw(screen *ebiten.Image) {
	m.drawMainTitle(screen)

	submenuTitleTextFace := *m.mainTitleTextFace
	submenuTitleTextFace.Size = 20

	switch m.state {
	case mainMenu:
		m.drawText(screen, &submenuTitleTextFace, true, localModeStr, multiplayerStr)
	case localMode:
		m.drawText(screen, &submenuTitleTextFace, true, onePlayerStr, twoPlayersStr, backStr)
	case multiplayerMode:
		m.drawText(screen, &submenuTitleTextFace, false, multiplayerFindOpponentStr)
	case levelSelection:
		m.drawText(screen, &submenuTitleTextFace, true, levelEasyStr, levelMediumStr, levelHardStr, backStr)
	}
}

func (m *Menu) drawMainTitle(screen *ebiten.Image) {
	positionX, _ := text.Measure(titleStr, m.mainTitleTextFace, 1)

	uiText := ui.Text{
		Value:    titleStr,
		FontFace: m.mainTitleTextFace,
		Position: geometry.Vector{
			X: (m.screenWidth - positionX) / 2,
			Y: 80,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)
}

func (m *Menu) drawText(screen *ebiten.Image, font *text.GoTextFace, drawBullet bool, values ...string) {
	var maxPositionX float64
	y := 250.0

	for _, val := range values {
		positionX, _ := text.Measure(val, font, 1)

		uiText := ui.Text{
			Value:    val,
			FontFace: font,
			Position: geometry.Vector{
				X: (m.screenWidth - positionX) / 2,
				Y: y,
			},
			Color: ui.DefaultColor,
		}

		uiText.Draw(screen)

		y += 50

		maxPositionX = math.Max(maxPositionX, positionX)
	}

	if !drawBullet {
		return
	}

	// draw selected option
	switch m.selectedOption {
	case 0:
		y = 255
	case 1:
		y = 305
	case 2:
		y = 355
	case 3:
		y = 405
	}

	vector.DrawFilledRect(
		screen,
		float32(m.screenWidth-maxPositionX)/2-30, float32(y),
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
