package menu

import (
	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	"github.com/hajimehoshi/ebiten/v2"
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
	font         *font.Font
	gameMode     GameMode
	level        level.Level
	readyToPlay  bool
	playerName   string
	screenHeight int
	screenWidth  int
	currentState state
	// states act as a cache to avoid creating the same state multiple times.
	states map[string]state
}

// New creates a new game menu.
func New(font *font.Font, screenWidth, screenHeight int) *Menu {
	menu := &Menu{
		font:         font,
		gameMode:     Undefined,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		states:       make(map[string]state),
	}

	menu.ChangeState(newMainMenuState(menu))

	return menu
}

// ChangeState changes the current state of the menu.
func (m *Menu) ChangeState(state state) {
	if st, ok := m.states[state.String()]; ok {
		m.currentState = st
		return
	}

	m.states[state.String()] = state
	m.currentState = state
}

// Update updates the menu.
func (m *Menu) Update() {
	m.currentState.Update()
}

// Draw draws the menu.
func (m *Menu) Draw(screen *ebiten.Image) {
	m.currentState.Draw(screen)
}

// IsReadyToPlay returns if the game is ready to play.
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

// PlayerName returns the given player name.
// This is only used in the multiplayer game mode.
func (m *Menu) PlayerName() string {
	return m.playerName
}
