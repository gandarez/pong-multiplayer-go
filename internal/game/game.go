package game

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/ai"
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	engineball "github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	engineplayer "github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	// ScreenWidth is the width of the screen.
	ScreenWidth = 640
	// ScreenHeight is the height of the screen.
	ScreenHeight = 480
)

// Game represents the game and implements the ebiten.Game interface.
type Game struct {
	assets *assets.Assets
	menu   *menu.Menu

	ball     *ball
	nextSide geometry.Side // it will be used to determine which side will start the game

	// players
	player1 *player
	player2 *player

	// scores
	score1 *score
	score2 *score

	// cpu
	cpu *ai.CPU
}

// New creates a new game.
func New(assets *assets.Assets) (*Game, error) {
	menu, err := menu.New(assets, ScreenWidth)
	if err != nil {
		return nil, fmt.Errorf("failed to create main menu: %w", err)
	}

	scoreTextFaceSource, err := assets.NewTextFaceSource("score")
	if err != nil {
		return nil, fmt.Errorf("failed to create score text face source: %w", err)
	}

	pongScoreFontFace := &text.GoTextFace{
		Source: scoreTextFaceSource,
		Size:   44,
	}

	p1 := engineplayer.New("Player 1", geometry.Left, ScreenWidth, ScreenHeight, 10)
	p2 := engineplayer.New("Player 2", geometry.Right, ScreenWidth, ScreenHeight, 10)

	score1AdjustmentPositionX, _ := text.Measure("0", pongScoreFontFace, 1)

	var nextPlayer geometry.Side
	if rand.Intn(2) == 0 { // nolint: gosec
		nextPlayer = geometry.Left
	} else {
		nextPlayer = geometry.Right
	}

	return &Game{
		assets: assets,
		ball: &ball{
			engineball.New(nextPlayer, ScreenWidth, ScreenHeight),
		},
		cpu:      ai.NewCPU(),
		menu:     menu,
		nextSide: nextPlayer,
		player1:  &player{p1},
		player2:  &player{p2},
		score1: &score{
			textFace: pongScoreFontFace,
			position: geometry.Vector{
				X: ScreenWidth/2 - 50 - score1AdjustmentPositionX,
				Y: 30,
			},
		},
		score2: &score{
			textFace: pongScoreFontFace,
			position: geometry.Vector{
				X: ScreenWidth/2 + 70,
				Y: 30,
			},
		},
	}, nil
}

// Update updates the game.
func (g *Game) Update() error {
	// if the game is not ready to play, udpate the menu
	if !g.menu.IsReadyToPlay() {
		g.menu.Update()
	}

	// update the ball
	g.ball.Update(g.player1.Bounds(), g.player2.Bounds())

	// update the players
	switch g.menu.GameMode() {
	case menu.OnePlayer:
		// guess ball position
		y := g.cpu.GuessBallPosition(
			g.ball.Bounds().Y,
			g.player1.Position().Y,
			g.player1.BouncerHeight(),
			ScreenHeight,
		)

		// ai player
		g.player1.SetPosition(y)
		// human player
		g.player2.Update(ebiten.KeyUp, ebiten.KeyDown)
	case menu.TwoPlayers:
		g.player1.Update(ebiten.KeyQ, ebiten.KeyA)
		g.player2.Update(ebiten.KeyUp, ebiten.KeyDown)
	}

	// check if the ball is out of the field and update the scores
	if out, side := g.ball.CheckGoal(); out {
		if side == geometry.Left {
			g.score2.value++
		} else {
			g.score1.value++
		}

		if g.nextSide == geometry.Left {
			g.nextSide = geometry.Right
		} else {
			g.nextSide = geometry.Left
		}

		g.ball = &ball{g.ball.Reset(g.nextSide)}
	}

	return nil
}

// Draw draws the game.
func (g *Game) Draw(screen *ebiten.Image) {
	// if the game is not ready to play, draw the menu
	if !g.menu.IsReadyToPlay() {
		g.menu.Draw(screen)

		return
	}

	// draw the field
	g.drawField(screen)

	// draw the ball
	g.ball.draw(screen)

	// draw the scores
	g.score1.draw(screen)
	g.score2.draw(screen)

	// draw the players
	g.player1.draw(screen)
	g.player2.draw(screen)
}

// Layout returns the screen width and height.
func (Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}
