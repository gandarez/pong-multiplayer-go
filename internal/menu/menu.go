package menu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	titleStr                   = "PONGO"
	localModeStr               = "Local Mode"
	multiplayerStr             = "Multiplayer"
	multiplayerFindOpponentStr = "Finding opponent.."
	onePlayerStr               = "One Player (Deactivated)"
	twoPlayersStr              = "Two Players"
	backStr                    = "Back"
)

// state is the state of the menu.
type state int

const (
	mainMenu state = iota
	localMode
	multiplayerMode
)

// gameMode is the game mode.
type gameMode int

const (
	undefined gameMode = iota
	onePlayer
	twoPlayers
	multiplayer
)

// Menu represents the game menu.
type Menu struct {
	mainTitleTextFace *text.GoTextFace
	gameMode          gameMode
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
		gameMode:          undefined,
		state:             mainMenu,
		screenWidth:       float64(screenWidth),
	}, nil
}

// Draw draws the menu on the screen.
func (m *Menu) Draw(screen *ebiten.Image) {
	mainTitleAdjustmentPositionX, _ := text.Measure(titleStr, m.mainTitleTextFace, 1)

	mainTitle := ui.Text{
		Value:    titleStr,
		FontFace: m.mainTitleTextFace,
		Position: geometry.Vector{
			X: (m.screenWidth - mainTitleAdjustmentPositionX) / 2,
			Y: 80,
		},
		Color: color.RGBA{200, 200, 200, 255},
	}

	mainTitle.Draw(screen)

	submenuTitleTextFace := *m.mainTitleTextFace
	submenuTitleTextFace.Size = 20

	switch m.state {
	case mainMenu:
		submenuAdjustmentPositionX, _ := text.Measure(localModeStr, &submenuTitleTextFace, 1)

		localModeTitleTextFace := ui.Text{
			Value:    localModeStr,
			FontFace: &submenuTitleTextFace,
			Position: geometry.Vector{
				X: (m.screenWidth - submenuAdjustmentPositionX) / 2,
				Y: 300,
			},
			Color: color.RGBA{200, 200, 200, 255},
		}

		localModeTitleTextFace.Draw(screen)

		submenuAdjustmentPositionX, _ = text.Measure(multiplayerStr, &submenuTitleTextFace, 1)

		multiplayerTitleTextFace := ui.Text{
			Value:    multiplayerStr,
			FontFace: &submenuTitleTextFace,
			Position: geometry.Vector{
				X: (m.screenWidth - submenuAdjustmentPositionX) / 2,
				Y: 350,
			},
			Color: color.RGBA{200, 200, 200, 255},
		}

		multiplayerTitleTextFace.Draw(screen)

		// draw selected option
		switch m.selectedOption {
		case 0:
			vector.DrawFilledRect(
				screen,
				float32(m.screenWidth-submenuAdjustmentPositionX)/2-30, 305,
				15, 15, color.RGBA{200, 200, 200, 255}, true,
			)
		case 1:
			vector.DrawFilledRect(
				screen,
				float32(m.screenWidth-submenuAdjustmentPositionX)/2-30, 355,
				15, 15, color.RGBA{200, 200, 200, 255}, true,
			)
		}
	case localMode:
		singlePlayerAdjustmentPositionX, _ := text.Measure(onePlayerStr, &submenuTitleTextFace, 1)

		localModeTitleTextFace := ui.Text{
			Value:    onePlayerStr,
			FontFace: &submenuTitleTextFace,
			Position: geometry.Vector{
				X: (m.screenWidth - singlePlayerAdjustmentPositionX) / 2,
				Y: 300,
			},
			Color: color.RGBA{200, 200, 200, 255},
		}

		localModeTitleTextFace.Draw(screen)

		twoPlayersAdjustmentPositionX, _ := text.Measure(twoPlayersStr, &submenuTitleTextFace, 1)

		multiplayerTitleTextFace := ui.Text{
			Value:    twoPlayersStr,
			FontFace: &submenuTitleTextFace,
			Position: geometry.Vector{
				X: (m.screenWidth - twoPlayersAdjustmentPositionX) / 2,
				Y: 350,
			},
			Color: color.RGBA{200, 200, 200, 255},
		}

		multiplayerTitleTextFace.Draw(screen)

		backMenuAdjustmentPositionX, _ := text.Measure(backStr, &submenuTitleTextFace, 1)

		backMenuTitleTextFace := ui.Text{
			Value:    backStr,
			FontFace: &submenuTitleTextFace,
			Position: geometry.Vector{
				X: (m.screenWidth - backMenuAdjustmentPositionX) / 2,
				Y: 400,
			},
			Color: color.RGBA{200, 200, 200, 255},
		}

		backMenuTitleTextFace.Draw(screen)

		// draw selected option
		switch m.selectedOption {
		case 0:
			vector.DrawFilledRect(
				screen,
				float32(m.screenWidth-singlePlayerAdjustmentPositionX)/2-30, 305,
				15, 15, color.RGBA{200, 200, 200, 255}, true,
			)
		case 1:
			vector.DrawFilledRect(
				screen,
				float32(m.screenWidth-singlePlayerAdjustmentPositionX)/2-30, 355,
				15, 15, color.RGBA{200, 200, 200, 255}, true,
			)
		case 2:
			vector.DrawFilledRect(
				screen,
				float32(m.screenWidth-singlePlayerAdjustmentPositionX)/2-30, 405,
				15, 15, color.RGBA{200, 200, 200, 255}, true,
			)
		}
	case multiplayerMode:
		submenuAdjustmentPositionX, _ := text.Measure(multiplayerFindOpponentStr, &submenuTitleTextFace, 1)

		finding := ui.Text{
			Value:    multiplayerFindOpponentStr,
			FontFace: &submenuTitleTextFace,
			Position: geometry.Vector{
				X: (m.screenWidth - submenuAdjustmentPositionX) / 2,
				Y: 300,
			},
			Color: color.RGBA{200, 200, 200, 255},
		}

		finding.Draw(screen)
	}
}

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
				// m.gameMode = onePlayer
				// m.readyToPlay = true
			case 1:
				// two players
				m.gameMode = twoPlayers
				m.readyToPlay = true
			case 2:
				m.state = mainMenu
				m.gameMode = undefined
				m.readyToPlay = false
				m.selectedOption = 0
			}
		}

		return
	}
}

func (m *Menu) IsReadyToPlay() bool {
	return m.readyToPlay
}

func (m *Menu) GameMode() gameMode {
	return m.gameMode
}
