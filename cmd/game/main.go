package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/game"
)

const title = "PONGO"

func main() {
	ebiten.SetWindowSize(int(game.ScreenWidth)*2, int(game.ScreenHeight)*2)
	ebiten.SetWindowTitle(title)
	ebiten.SetRunnableOnUnfocused(true)

	// load all assets
	assets, err := assets.Load()
	if err != nil {
		slog.Error("failed to load assets", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gameInstance, err := game.New(ctx, cancel, assets)
	if err != nil {
		slog.Error("failed to create game", slog.Any("error", err))
		os.Exit(1) // nolint:gocritic
	}

	// run the game and lock the main goroutine
	if err := ebiten.RunGame(gameInstance); err != nil {
		slog.Error("failed to run game", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("exiting the game")
}
