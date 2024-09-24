package game

import (
	"context"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/network"
)

const (
	// maxScore is the maximum score to win the game.
	maxScore = 1
	// ScreenWidth is the width of the screen.
	ScreenWidth = 640
	// ScreenHeight is the height of the screen.
	ScreenHeight = 480
	// fieldBorderWidth is the width of the field border.
	fieldBorderWidth = 10
)

// Game represents the main game object.
type Game struct {
	cancel context.CancelFunc
	ctx    context.Context
	font   *font.Font
	menu   *menu.Menu

	currentState state

	// Shared resources
	assets        *assets.Assets
	networkClient *network.Client
}

// New creates a new game instance.
func New(ctx context.Context, cancel context.CancelFunc, assets *assets.Assets) (*Game, error) {
	font := font.New(assets)
	gameMenu := menu.New(font, ScreenWidth)

	game := &Game{
		cancel: cancel,
		ctx:    ctx,
		font:   font,
		menu:   gameMenu,
		assets: assets,
	}

	// Set the initial state to MainMenuState
	game.currentState = newMainMenuState(game)

	return game, nil
}

// Update delegates the update logic to the current game state.
func (g *Game) Update() error {
	if err := g.currentState.update(); err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}

	return nil
}

// Draw delegates the drawing logic to the current game state.
func (g *Game) Draw(screen *ebiten.Image) {
	g.currentState.draw(screen)
}

// Layout returns the game's logical screen dimensions.
func (*Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// changeState allows switching between different game states.
func (g *Game) changeState(state state) {
	g.currentState = state
}

// exit gracefully exits the game.
func (g *Game) exit() error {
	if g.networkClient != nil {
		g.networkClient.Close()
	}

	if g.cancel != nil {
		g.cancel()
	}

	// gracefully exit the game
	return ebiten.Termination
}
