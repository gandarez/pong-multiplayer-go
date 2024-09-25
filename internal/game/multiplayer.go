package game

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/internal/font"
	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
)

// multiplayerState represents the multiplayer game state.
type multiplayerState struct {
	ball                           ball.Ball
	player1                        player.Player
	player2                        player.Player
	score1                         *score
	score2                         *score
	networkGameCh                  chan network.GameState
	p1NamePosition, p2NamePosition geometry.Vector
	*baseState
}

// newMultiplayerState creates a new multiplayerState.
func newMultiplayerState(game *Game, ready network.ReadyMessage) *multiplayerState {
	base := newBasePlayingState(game, game.menu.Level())

	// initialize players with names from gameState
	player1 := player.NewNetwork(ready.Name, ready.Side, ScreenWidth, ScreenHeight)
	player2 := player.NewNetwork(ready.OpponentName, ready.OpponentSide, ScreenWidth, ScreenHeight)
	ball := ball.NewNetwork()
	score1 := newScore1(base.game.font)
	score2 := newScore2(base.game.font)

	// calculate player name position
	p1NamePosition, p2NamePosition := calculatePlayerNamePosition(*game.font, player1.Name(), player2.Name())

	networkGameCh := make(chan network.GameState)

	go func() {
		if err := game.networkClient.ReceiveGameState(networkGameCh); err != nil {
			slog.Error("failed to receive game state", slog.Any("error", err))
		}
	}()

	return &multiplayerState{
		baseState:      base,
		ball:           ball,
		player1:        player1,
		player2:        player2,
		score1:         score1,
		score2:         score2,
		networkGameCh:  networkGameCh,
		p1NamePosition: p1NamePosition,
		p2NamePosition: p2NamePosition,
	}
}

func (s *multiplayerState) update() error {
	// update common elements
	s.baseState.update()

	if s.gamePaused {
		return nil
	}

	up := ebiten.IsKeyPressed(ebiten.KeyUp)
	down := ebiten.IsKeyPressed(ebiten.KeyDown)

	if up || down {
		// send input to server
		if err := s.game.networkClient.SendPlayerInput(network.PlayerInput{
			Up:   up,
			Down: down,
		}); err != nil {
			slog.Error("failed to send player input", slog.Any("error", err))
		}
	}

	// receive game state from server and update local game state
	gameState := <-s.networkGameCh

	// update ball and players positions
	s.updateBallTrail(s.ball)
	s.ball.SetPosition(gameState.Ball.Position)
	s.ball.SetAngle(gameState.Ball.Angle)
	s.ball.SetBounces(gameState.Ball.Bounces)

	// update player positions and scores
	s.updatePlayerPositions(gameState)
	s.updateScores(gameState)

	// update ping
	s.pingCurrentPlayer = gameState.CurrentPlayer.Ping
	s.pingOpponent = gameState.OpponentPlayer.Ping

	// check for winner
	if gameState.CurrentPlayer.Winner || gameState.OpponentPlayer.Winner {
		winner := s.player1
		if gameState.OpponentPlayer.Winner {
			winner = s.player2
		}

		s.game.networkClient.Close()

		s.game.changeState(newWinnerState(s.game, winner.Name(), s))
	}

	return nil
}

func (s *multiplayerState) updatePlayerPositions(gameState network.GameState) {
	if s.player1.Side() == gameState.CurrentPlayer.Side {
		s.player1.SetPosition(gameState.CurrentPlayer.PositionY)
		s.player2.SetPosition(gameState.OpponentPlayer.PositionY)
	} else {
		s.player1.SetPosition(gameState.OpponentPlayer.PositionY)
		s.player2.SetPosition(gameState.CurrentPlayer.PositionY)
	}
}

func (s *multiplayerState) updateScores(gameState network.GameState) {
	if s.player1.Side() == gameState.CurrentPlayer.Side {
		s.score1.value = gameState.CurrentPlayer.Score
		s.score2.value = gameState.OpponentPlayer.Score
	} else {
		s.score1.value = gameState.OpponentPlayer.Score
		s.score2.value = gameState.CurrentPlayer.Score
	}
}

func calculatePlayerNamePosition(font font.Font, p1name, p2name string) (geometry.Vector, geometry.Vector) {
	playerNameTextFace, _ := font.Face("ui", 20)
	scoreTextFace, _ := font.Face("score", 44)

	p1NameWidth, _ := text.Measure(fmt.Sprintf("%10s", p1name), playerNameTextFace, 1)
	p2NameWidth, _ := text.Measure(fmt.Sprintf("%10s", p2name), playerNameTextFace, 1)

	_, scoreHeight := text.Measure("0", scoreTextFace, 1)

	p1NamePosition := geometry.Vector{
		X: ScreenWidth/2 - 10 - p1NameWidth,
		Y: scoreHeight + 50,
	}

	p2NamePosition := geometry.Vector{
		X: ScreenWidth/2 + 10 + (p2NameWidth / 2),
		Y: scoreHeight + 50,
	}

	return p1NamePosition, p2NamePosition
}

func (s *multiplayerState) draw(screen *ebiten.Image) {
	// draw common elements
	s.baseState.draw(screen)

	// draw players, ball, and scores
	drawPlayer(s.player1.Position(), s.player1.BouncerWidth(), s.player1.BouncerHeight(), screen)
	drawPlayer(s.player2.Position(), s.player2.BouncerWidth(), s.player2.BouncerHeight(), screen)
	drawBall(screen, s.ball.Position(), s.ball.Width(), s.ballTrail)
	s.score1.draw(screen)
	s.score2.draw(screen)

	// draw player1 name
	if err := drawPlayerName(s.player1.Name(), s.p1NamePosition, screen, s.game.font); err != nil {
		slog.Error("failed to draw player name", slog.Any("error", err))
		panic(err)
	}

	// draw player2 name
	if err := drawPlayerName(s.player2.Name(), s.p2NamePosition, screen, s.game.font); err != nil {
		slog.Error("failed to draw player name", slog.Any("error", err))
		panic(err)
	}

	// draw metric
	s.metric.DrawNetworkInfo(screen, s.pingCurrentPlayer, s.pingOpponent)
}

func (s *multiplayerState) getBall() ball.Ball {
	return s.ball
}

func (*multiplayerState) canPause() bool {
	return false
}
