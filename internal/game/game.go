package game

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
	"github.com/gandarez/pong-multiplayer-go/internal/ai"
	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/network"
	metric "github.com/gandarez/pong-multiplayer-go/internal/stat"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	engineball "github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
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
	// connecting is used when the game is connecting to connect to the party.
	connecting
	// playing is used when the game is being played.
	playing
	// ended is used when the game is over.
	ended
)

// Game represents the game and implements the ebiten.Game interface.
type Game struct {
	cancel     context.CancelFunc
	ctx        context.Context
	font       *font.Font
	menu       *menu.Menu
	metric     *metric.Metric
	showMetric bool
	state      state

	ball     *ball
	nextSide geometry.Side // it will be used to determine which side will start the game

	// players
	player1 *player
	player2 *player

	// scores
	score1 *score
	score2 *score

	// multiplayer
	networkClient      *network.Client
	networkGameStateCh chan network.GameState
	pingCurrentPlayer  int64
	pingOpponent       int64

	// ready is used when the game is ready to play to create some objects only once
	ready sync.Once
}

// New creates a new game.
func New(ctx context.Context, cancel context.CancelFunc, assets *assets.Assets) (*Game, error) {
	font := font.New(assets)

	// create the metric
	metric, err := metric.New(font, ScreenWidth)
	if err != nil {
		slog.Error("failed to create metric", slog.Any("error", err))
	}

	return &Game{
		cancel:             cancel,
		ctx:                ctx,
		font:               font,
		menu:               menu.New(font, ScreenWidth),
		metric:             metric,
		networkGameStateCh: make(chan network.GameState),
		ready:              sync.Once{},
		state:              notReady,
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

		// it means to exit the game
		if g.menu.GameMode() == menu.Undefined {
			if runtime.GOOS != "js" {
				return ebiten.Termination
			}

			os.Exit(0)
		}

		// if game mode is multiplayer, then change the state to connecting and wait for the connection
		if g.menu.GameMode() == menu.Multiplayer {
			g.state = connecting
			return nil
		}

		// if the game is ready to play, change the state to playing
		g.state = playing
	case connecting:
		if g.networkClient == nil {
			g.networkClient = network.NewClient(g.ctx, g.cancel, network.BaseURL)
			if err := g.networkClient.Connect(); err != nil {
				return fmt.Errorf("failed to connect to the server: %w", err)
			}

			// TODO: get player name from input
			if err := g.networkClient.SendPlayerInfo(network.GameInfo{
				PlayerName:       "Player 1",
				Level:            int(g.menu.Level()),
				ScreenWidth:      ScreenWidth,
				ScreenHeight:     ScreenHeight,
				MaxScore:         maxScore,
				FieldBorderWidth: int(fieldBorderWidth),
			}); err != nil {
				return fmt.Errorf("failed to send player info: %w", err)
			}

			go g.networkClient.ReceiveGameState(g.networkGameStateCh)

			go func() {
				// wait for the connection to be established
				<-g.networkGameStateCh

				// TODO: check if start returns an error
				_ = g.start()
				g.state = playing
			}()
		}
	case playing:
		if g.menu.GameMode() == menu.Multiplayer {
			gameState := <-g.networkGameStateCh

			// update the game state
			if gameState.CurrentPlayer.Side == geometry.Left {
				g.player1.SetPosition(gameState.CurrentPlayer.PositionY)
			} else {
				g.player2.SetPosition(gameState.CurrentPlayer.PositionY)
			}

			if gameState.OpponentPlayer.Side == geometry.Left {
				g.player1.SetPosition(gameState.OpponentPlayer.PositionY)
			} else {
				g.player2.SetPosition(gameState.OpponentPlayer.PositionY)
			}

			g.ball.SetAngle(gameState.Ball.Angle)
			g.ball.SetBounces(gameState.Ball.Bounces)
			g.ball.SetPosition(gameState.Ball.Position)

			g.score1.value = gameState.CurrentPlayer.Score
			g.score2.value = gameState.OpponentPlayer.Score

			g.pingCurrentPlayer = gameState.CurrentPlayer.Ping
			g.pingOpponent = gameState.OpponentPlayer.Ping
		}

		// call the update to do the game logic
		if err := g.update(); err != nil {
			return fmt.Errorf("failed to update the game: %w", err)
		}

		// if ESC is pressed, then show/hide metrics
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			g.showMetric = !g.showMetric
		}
	case ended:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
	}

	return nil
}

func (g *Game) update() error {
	if g.menu.GameMode() != menu.Multiplayer {
		// update the ball
		g.ball.Update(g.player1.Bounds(), g.player2.Bounds())
	}

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
		g.player2.Update(engineplayer.Input{
			Up:   ebiten.IsKeyPressed(ebiten.KeyUp),
			Down: ebiten.IsKeyPressed(ebiten.KeyDown),
		})
	case menu.TwoPlayers:
		g.player1.Update(engineplayer.Input{
			Up:   ebiten.IsKeyPressed(ebiten.KeyQ),
			Down: ebiten.IsKeyPressed(ebiten.KeyA),
		})
		g.player2.Update(engineplayer.Input{
			Up:   ebiten.IsKeyPressed(ebiten.KeyUp),
			Down: ebiten.IsKeyPressed(ebiten.KeyDown),
		})
	case menu.Multiplayer:
		up := ebiten.IsKeyPressed(ebiten.KeyUp)
		down := ebiten.IsKeyPressed(ebiten.KeyDown)

		// only send the player input if the user is pressing the keys
		if up || down {
			if err := g.networkClient.SendPlayerInput(network.PlayerInput{
				Up:   up,
				Down: down,
			}); err != nil {
				slog.Error("failed to send player input", slog.Any("error", err))
			}
		}
	}

	// eraly return if the game mode is multiplayer
	if g.menu.GameMode() == menu.Multiplayer {
		return nil
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
		if g.networkClient != nil {
			g.networkClient.Close()
		}

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
			ui.DrawSplash(screen, g.font, ScreenWidth)
			g.menu.Draw(screen)
		}
	case connecting:
		ui.DrawSplash(screen, g.font, ScreenWidth)
		ui.DrawWaitingConnection(screen, g.font, ScreenWidth)
	case playing:
		g.draw(screen)
	case ended:
		g.draw(screen) // draw the current state of the game when the game is over
		g.drawWinner(screen)
	}
}

func (g *Game) draw(screen *ebiten.Image) {
	// initialize the game. It's thread-safe.
	if err := g.start(); err != nil {
		// panic is not the best way to handle this error, but Draw() does not return an error
		panic(fmt.Errorf("failed to start the game: %w", err))
	}

	// draw the field
	g.drawField(screen)

	// draw the metric
	g.tryDrawMetric(screen)

	// draw the ball
	g.ball.draw(screen)

	// draw the scores
	g.score1.draw(screen)
	g.score2.draw(screen)

	// draw the players
	g.player1.draw(screen)
	g.player2.draw(screen)

	// draw player names when multiplayer mode
	if g.menu.GameMode() == menu.Multiplayer {
		g.player1.drawName(screen, g.font)
		g.player2.drawName(screen, g.font)
	}
}

func (g *Game) tryDrawMetric(screen *ebiten.Image) {
	if g.metric == nil {
		slog.Warn("metric is nil")
	}

	if !g.showMetric {
		return
	}

	g.metric.Draw(screen, g.ball.Bounces(), g.ball.Angle(), g.menu.Level())

	if g.menu.GameMode() == menu.Multiplayer {
		g.metric.DrawNetworkInfo(screen, g.pingCurrentPlayer, g.pingOpponent)
	}
}

// Layout returns the screen width and height.
func (*Game) Layout(_, _ int) (int, int) {
	return int(ScreenWidth), int(ScreenHeight)
}

// gameEnded checks if the game is over.
func (g *Game) gameEnded() bool {
	return g.score1.value == maxScore || g.score2.value == maxScore
}

func (g *Game) start() (errstart error) {
	g.ready.Do(func() {
		scoreTextFace, err := g.font.Face("score", 44)
		if err != nil {
			errstart = fmt.Errorf("failed to create score text face: %w", err)
			return
		}

		playerNameTextFace, err := g.font.Face("ui", 20)
		if err != nil {
			errstart = fmt.Errorf("failed to create player name text face: %w", err)
			return
		}

		var nextPlayer geometry.Side
		if rand.Intn(2) == 0 { // nolint:gosec
			nextPlayer = geometry.Left
		} else {
			nextPlayer = geometry.Right
		}

		g.nextSide = nextPlayer

		scoreWidth, scoreHeight := text.Measure("0", scoreTextFace, 1)

		var (
			b                              engineball.Ball
			p1, p2                         engineplayer.Player
			p1NamePosition, p2NamePosition geometry.Vector // only for multiplayer mode
		)

		switch g.menu.GameMode() {
		case menu.OnePlayer:
			p1 = engineplayer.NewLocal("CPU", geometry.Left, ScreenWidth, ScreenHeight, fieldBorderWidth)
			p2 = engineplayer.NewLocal("Player", geometry.Right, ScreenWidth, ScreenHeight, fieldBorderWidth)

			b = engineball.NewLocal(nextPlayer, ScreenWidth, ScreenHeight, g.menu.Level())
		case menu.TwoPlayers:
			p1 = engineplayer.NewLocal("Player 1", geometry.Left, ScreenWidth, ScreenHeight, fieldBorderWidth)
			p2 = engineplayer.NewLocal("Player 2", geometry.Right, ScreenWidth, ScreenHeight, fieldBorderWidth)

			b = engineball.NewLocal(nextPlayer, ScreenWidth, ScreenHeight, g.menu.Level())
		case menu.Multiplayer:
			p1 = engineplayer.NewNetwork("Player 1", geometry.Left, ScreenWidth, ScreenHeight)
			p2 = engineplayer.NewNetwork("Player 2", geometry.Right, ScreenWidth, ScreenHeight)

			p1NameWidth, _ := text.Measure(p1.Name(), playerNameTextFace, 1)
			p2NameWidth, _ := text.Measure(p2.Name(), playerNameTextFace, 1)

			p1NamePosition = geometry.Vector{
				X: ScreenWidth/2 - 10 - p1NameWidth,
				Y: scoreHeight + 40,
			}

			p2NamePosition = geometry.Vector{
				X: ScreenWidth/2 + (p2NameWidth / 2),
				Y: scoreHeight + 40,
			}

			b = engineball.NewNetwork()
		}

		g.player1 = &player{p1NamePosition, p1}
		g.player2 = &player{p2NamePosition, p2}
		g.ball = &ball{
			b,
		}

		g.score1 = &score{
			textFace: scoreTextFace,
			position: geometry.Vector{
				X: ScreenWidth/2 - 50 - scoreWidth,
				Y: 30,
			},
		}
		g.score2 = &score{
			textFace: scoreTextFace,
			position: geometry.Vector{
				X: ScreenWidth/2 + 70,
				Y: 30,
			},
		}
	})

	return //nolint:revive
}
