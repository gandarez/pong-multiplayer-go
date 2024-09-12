package game

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/ai"
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	engineball "github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/level"
	engineplayer "github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

const (
	// maxScore is the maximum score to win the game.
	maxScore = 10
	// ScreenWidth is the width of the screen.
	ScreenWidth = 640
	// ScreenHeight is the height of the screen.
	ScreenHeight = 480
)

// state represents the state of the game.
type state int

const (
	// notReady is used to show the main menu.
	notReady state = iota
	// playing is used when the game is being played.
	playing
	// ended is used when the game is over.
	ended
)

// Game represents the game and implements the ebiten.Game interface.
type Game struct {
	assets *assets.Assets
	menu   *menu.Menu
	state  state

	ball     *ball
	nextSide geometry.Side // it will be used to determine which side will start the game

	// players
	player1 *player
	player2 *player

	// scores
	score1 *score
	score2 *score

	// ready is used when the game is ready to play to create some objects only once
	ready sync.Once
}

// New creates a new game.
func New(assets *assets.Assets) (*Game, error) {
	menu, err := menu.New(assets, ScreenWidth)
	if err != nil {
		return nil, fmt.Errorf("failed to create main menu: %w", err)
	}

	return &Game{
		assets: assets,
		menu:   menu,
		ready:  sync.Once{},
		state:  notReady,
	}, nil
}

func (g *Game) reset() {
	g.score1.value = 0
	g.score2.value = 0

	g.player1.Reset()
	g.player2.Reset()

	g.ball = &ball{g.ball.Reset(g.nextSide)}

	g.state = playing
}

// Update updates the game.
func (g *Game) Update() error {
	switch g.state {
	case notReady:
		// if the game is not ready to play, udpate the menu
		if !g.menu.IsReadyToPlay() {
			g.menu.Update()

			return nil
		}

		// if the game is ready to play, change the state to playing
		g.state = playing
	case playing:
		g.update()
	case ended:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
	}

	return nil
}

func (g *Game) update() error {
	// update the ball
	g.ball.Update(g.player1.Bounds(), g.player2.Bounds())

	// update the players
	switch g.menu.GameMode() {
	case menu.OnePlayer:
		// guess ball position
		y := ai.GuessBallPosition(
			g.ball.Bounds().Y,
			g.player1.Position().Y,
			g.player1.BouncerHeight(),
			ScreenHeight,
		)

		// AI player
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

	// game has ended?
	if g.gameEnded() {
		g.state = ended
	}

	return nil
}

// Draw draws the game.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case notReady:
		// if the game is not ready to play, draw the menu
		if !g.menu.IsReadyToPlay() {
			g.menu.Draw(screen)
		}
	case playing:
		g.draw(screen)
	case ended:
		g.draw(screen) // draw the state of the game
		g.drawWinner(screen)
	}
}

func (g *Game) draw(screen *ebiten.Image) {
	// initialize the game. It's thread-safe.
	if err := g.start(g.menu.Level()); err != nil {
		// panic is not the best way to handle this error, but Draw() does not return an error
		panic(fmt.Errorf("failed to start the game: %w", err))
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
func (*Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// gameEnded checks if the game is over.
func (g *Game) gameEnded() bool {
	return g.score1.value == maxScore || g.score2.value == maxScore
}

func (g *Game) drawWinner(screen *ebiten.Image) {
	var winner string
	if g.score1.value == maxScore {
		winner = g.player1.Name()
	} else {
		winner = g.player2.Name()
	}

	winnerFaceSource, err := g.assets.NewTextFaceSource("ui")
	if err != nil {
		// panic is not the best way to handle this error, but Draw() does not return an error
		panic(fmt.Errorf("failed to create winner text face source: %w", err))
	}

	font := &text.GoTextFace{
		Source: winnerFaceSource,
		Size:   40,
	}

	winnerText := fmt.Sprintf("%s won", winner)

	positionX, _ := text.Measure(winnerText, font, 1)

	uiText := ui.Text{
		Value:    fmt.Sprintf("%s won", winner),
		FontFace: font,
		Position: geometry.Vector{
			X: (ScreenWidth - positionX) / 2,
			Y: 200,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)

	font.Size = 30

	positionX, _ = text.Measure("Press Enter to play again", font, 1)

	uiText = ui.Text{
		Value:    "Press Enter to play again",
		FontFace: font,
		Position: geometry.Vector{
			X: (ScreenWidth - positionX) / 2,
			Y: 300,
		},
		Color: ui.DefaultColor,
	}

	uiText.Draw(screen)
}

func (g *Game) start(lvl level.Level) (errstart error) {
	g.ready.Do(func() {
		scoreTextFaceSource, err := g.assets.NewTextFaceSource("score")
		if err != nil {
			errstart = fmt.Errorf("failed to create score text face source: %w", err)
			return
		}

		p1 := engineplayer.New("Player 1", geometry.Left, ScreenWidth, ScreenHeight, 10)
		p2 := engineplayer.New("Player 2", geometry.Right, ScreenWidth, ScreenHeight, 10)

		g.player1 = &player{p1}
		g.player2 = &player{p2}

		pongScoreFontFace := &text.GoTextFace{
			Source: scoreTextFaceSource,
			Size:   44,
		}

		score1AdjustmentPositionX, _ := text.Measure("0", pongScoreFontFace, 1)

		g.score1 = &score{
			textFace: pongScoreFontFace,
			position: geometry.Vector{
				X: ScreenWidth/2 - 50 - score1AdjustmentPositionX,
				Y: 30,
			},
		}
		g.score2 = &score{
			textFace: pongScoreFontFace,
			position: geometry.Vector{
				X: ScreenWidth/2 + 70,
				Y: 30,
			},
		}

		var nextPlayer geometry.Side
		if rand.Intn(2) == 0 { // nolint: gosec
			nextPlayer = geometry.Left
		} else {
			nextPlayer = geometry.Right
		}

		g.nextSide = nextPlayer

		g.ball = &ball{
			engineball.New(nextPlayer, ScreenWidth, ScreenHeight, lvl),
		}
	})

	return
}
