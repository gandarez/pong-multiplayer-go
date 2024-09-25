package game

import (
	"log/slog"

	"github.com/gandarez/pong-multiplayer-go/internal/menu"
	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/ball"
	"github.com/gandarez/pong-multiplayer-go/pkg/engine/player"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type spectatorState struct {
	ball           ball.Ball
	player1        player.Player
	player2        player.Player
	score1         *score
	score2         *score
	gameStateCh    chan network.GameState
	sessionID      string
	p1NamePosition geometry.Vector
	p2NamePosition geometry.Vector
	*baseState
}

func newSpectatorState(game *Game) *spectatorState {
	base := newBasePlayingState(game, game.menu.Level())

	// initialize players and ball
	player1 := player.NewNetwork("", geometry.Left, ScreenWidth, ScreenHeight)
	player2 := player.NewNetwork("", geometry.Right, ScreenWidth, ScreenHeight)
	ball := ball.NewNetwork()
	score1 := newScore1(base.game.font) // left
	score2 := newScore2(base.game.font) // right

	p1NamePosition, p2NamePosition := calculatePlayerNamePosition(*game.font, player1.Name(), player2.Name(), player1.Side())

	gameStateCh := make(chan network.GameState)

	state := &spectatorState{
		baseState:      base,
		ball:           ball,
		player1:        player1,
		player2:        player2,
		score1:         score1,
		score2:         score2,
		gameStateCh:    gameStateCh,
		sessionID:      game.menu.SessionID,
		p1NamePosition: p1NamePosition,
		p2NamePosition: p2NamePosition,
	}

	state.connectAsSpectator()

	return state
}

func (s *spectatorState) update() error {
	// handle ESC key to go back to main menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.game.networkClient.Close()
		s.game.menu = menu.New(s.game.font, ScreenWidth, ScreenHeight)
		s.game.changeState(newMainMenuState(s.game))

		return nil
	}

	gameState := <-s.gameStateCh

	s.updateGameState(gameState)

	return nil
}

func (s *spectatorState) updateGameState(gameState network.GameState) {
	// update ball
	s.updateBallTrail(s.ball)
	s.ball.SetPosition(gameState.Ball.Position)
	s.ball.SetAngle(gameState.Ball.Angle)
	s.ball.SetBounces(gameState.Ball.Bounces)

	// update players
	s.player1.SetPosition(gameState.CurrentPlayer.PositionY)
	s.player2.SetPosition(gameState.OpponentPlayer.PositionY)

	// update player names if they have changed
	if s.player1.Name() != gameState.CurrentPlayer.Name {
		s.player1.SetName(gameState.CurrentPlayer.Name)
		s.updatePlayerNamePosition()
	}

	if s.player2.Name() != gameState.OpponentPlayer.Name {
		s.player2.SetName(gameState.OpponentPlayer.Name)
		s.updatePlayerNamePosition()
	}

	// update scores
	s.score1.value = gameState.CurrentPlayer.Score
	s.score2.value = gameState.OpponentPlayer.Score

	// check winner
	if gameState.CurrentPlayer.Winner || gameState.OpponentPlayer.Winner {
		winnerName := gameState.CurrentPlayer.Name
		if gameState.OpponentPlayer.Winner {
			winnerName = gameState.OpponentPlayer.Name
		}

		// close network connection
		s.game.networkClient.Close()

		// change state to winner screen
		s.game.changeState(newWinnerState(s.game, winnerName, s))
	}
}

func (s *spectatorState) draw(screen *ebiten.Image) {
	// draw common elements
	s.baseState.draw(screen)

	// draw players, ball, and scores
	drawPlayer(s.player1.Position(), s.player1.BouncerWidth(), s.player1.BouncerHeight(), screen)
	drawPlayer(s.player2.Position(), s.player2.BouncerWidth(), s.player2.BouncerHeight(), screen)
	drawBall(screen, s.ball.Position(), s.ball.Width(), s.ballTrail)
	s.score1.draw(screen)
	s.score2.draw(screen)

	// draw player names
	if err := drawPlayerName(s.player1.Name(), s.p1NamePosition, screen, s.game.font); err != nil {
		slog.Error("failed to draw player name", slog.Any("error", err))
	}

	if err := drawPlayerName(s.player2.Name(), s.p2NamePosition, screen, s.game.font); err != nil {
		slog.Error("failed to draw player name", slog.Any("error", err))
	}
}

func (s *spectatorState) getBall() ball.Ball {
	return s.ball
}

func (*spectatorState) canPause() bool {
	return false
}

func (s *spectatorState) connectAsSpectator() {
	s.game.networkClient = network.NewSpectatorClient(s.game.ctx, s.game.cancel)
	if err := s.game.networkClient.ConnectAsSpectator(s.sessionID); err != nil {
		slog.Error("failed to connect as spectator", slog.Any("error", err))
		s.game.changeState(newMainMenuState(s.game))

		return
	}

	go func() {
		if err := s.game.networkClient.ReceiveGameState(s.gameStateCh); err != nil {
			slog.Error("failed to receive game state", slog.Any("error", err))
		}
	}()
}

func (s *spectatorState) updatePlayerNamePosition() {
	s.p1NamePosition, s.p2NamePosition = calculatePlayerNamePosition(*s.game.font, s.player1.Name(), s.player2.Name(), s.player1.Side())
}
