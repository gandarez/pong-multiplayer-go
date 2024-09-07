package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/game"
)

const title = "PONGO"

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle(title)

	// load all assets
	assets, err := assets.Load()
	if err != nil {
		log.Fatalf("failed to load assets: %s", err)
	}

	game, err := game.New(assets)
	if err != nil {
		log.Fatalf("failed to create game: %s", err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
