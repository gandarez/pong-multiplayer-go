package game

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/network"
)

const (
	ScreenWidth      = 640
	ScreenHeight     = 480
	maxScore         = 10
	fieldBorderWidth = 10
)

var fieldColor = color.RGBA{200, 200, 200, 255}

// Game represents the main game object.
type Game struct {
	cancel context.CancelFunc
	ctx    context.Context
	font   *font.Font
	menu   *menu.Menu

	currentState GameState

	// Shared resources
	assets        *assets.Assets
	networkClient *network.Client
	networkGameCh chan network.GameState
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
	game.currentState = NewMainMenuState(game)

	return game, nil
}

// Update delegates the update logic to the current game state.
func (g *Game) Update() error {
	if err := g.currentState.Update(); err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}
	return nil
}

// Draw delegates the drawing logic to the current game state.
func (g *Game) Draw(screen *ebiten.Image) {
	g.currentState.Draw(screen)
}

// Layout returns the game's logical screen dimensions.
func (*Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// ChangeState allows switching between different game states.
func (g *Game) ChangeState(state GameState) {
	g.currentState = state
}

// Exit gracefully exits the game.
func (g *Game) Exit() {
	if g.networkClient != nil {
		g.networkClient.Close()
	}
	if g.cancel != nil {
		g.cancel()
	}
	if os.Getenv("EBITEN_RUN_IN_MAIN_THREAD") != "1" {
		os.Exit(0)
	} else {
		// For WebAssembly
		slog.Error("Exiting game")
	}
}

// drawField draws the game field, common to all game modes.
func (g *Game) drawField(screen *ebiten.Image) {
	// Draw field limits (top and bottom borders)
	for x := 0.; x <= ScreenWidth; x++ {
		for y := 0.; y <= fieldBorderWidth; y++ {
			screen.Set(int(x), int(y), fieldColor)
			screen.Set(int(x), int(y+ScreenHeight-fieldBorderWidth), fieldColor)
		}
	}

	// Draw center line (dashed)
	for squareCount, y := 0, 15.; squareCount < 30; squareCount++ {
		for w := 0.; w < 7.; w++ {
			for h := 0.; h < 7.; h++ {
				screen.Set(int(ScreenWidth/2-5+w), int(h+y), fieldColor)
			}
		}
		y += 17.
	}
}
