package menu

import (
	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
)

type GameMode int

const (
	Undefined GameMode = iota
	OnePlayer
	TwoPlayers
	Multiplayer
)

type Menu struct {
	font         *font.Font
	gameMode     GameMode
	level        level.Level
	readyToPlay  bool
	playerName   string
	screenWidth  int
	currentState State
}

func New(font *font.Font, screenWidth int) *Menu {
	menu := &Menu{
		font:        font,
		gameMode:    Undefined,
		playerName:  "",
		screenWidth: screenWidth,
	}
	menu.ChangeState(NewMainMenuState(menu))
	return menu
}

func (m *Menu) ChangeState(state State) {
	m.currentState = state
}

func (m *Menu) Update() {
	m.currentState.Update()
}

func (m *Menu) Draw(screen Screen) {
	m.currentState.Draw(screen)
}

func (m *Menu) IsReadyToPlay() bool {
	return m.readyToPlay
}

func (m *Menu) GameMode() GameMode {
	return m.gameMode
}

func (m *Menu) Level() level.Level {
	return m.level
}

func (m *Menu) PlayerName() string {
	return m.playerName
}
