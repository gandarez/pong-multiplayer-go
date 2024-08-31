package engine

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/geometry"
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	assets *assets.Assets
	menu   *menu.Menu

	ball     *ball
	nextSide string // it will be used to determine which side will start the game

	// players
	player1 *player
	player2 *player

	// scores
	score1 *score
	score2 *score
}

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

	p1, p2 := newPlayers()

	score1AdjustmentPositionX, _ := text.Measure("0", pongScoreFontFace, 1)

	var nextPlayer string
	if rand.Intn(2) == 0 {
		nextPlayer = player1Name
	} else {
		nextPlayer = player2Name
	}

	return &Game{
		assets:   assets,
		menu:     menu,
		ball:     newBall(nextPlayer),
		nextSide: nextPlayer,
		player1:  p1,
		player2:  p2,
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

func (g *Game) Update() error {
	// if the game is not ready to play, udpate the menu
	if !g.menu.IsReadyToPlay() {
		g.menu.Update()
	}

	// update the players
	g.player1.update(ebiten.KeyQ, ebiten.KeyA)
	g.player2.update(ebiten.KeyUp, ebiten.KeyDown)

	// update the ball
	g.ball.update()
	g.ball.bounce(g.player1.bounds(), g.player2.bounds())

	// check if the ball is out of the field and update the scores
	if out, side := g.ball.checkGoal(); out {
		if side == geometry.Left {
			g.score2.value++
		} else {
			g.score1.value++
		}

		if g.nextSide == player1Name {
			g.nextSide = player2Name
		} else {
			g.nextSide = player1Name
		}

		g.ball = newBall(g.nextSide)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
