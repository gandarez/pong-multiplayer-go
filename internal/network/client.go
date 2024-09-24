package network

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	// BaseURL is the base URL of the server.
	BaseURL = "game.go-go.dev"

	writeTimeout = 10 * time.Second
	readTimeout  = 60 * time.Second
)

// Client is a client that connects to the server using a websocket connection.
type Client struct {
	conn      *websocket.Conn
	serverURL string
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewClient creates a new client.
func NewClient(ctx context.Context, cancel context.CancelFunc, serverURL string) *Client {
	return &Client{
		serverURL: serverURL,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Connect connects to the server using a websocket connection.
func (c *Client) Connect() error {
	u := fmt.Sprintf("wss://%s/multiplayer", c.serverURL)

	ctx, cancel := context.WithTimeout(c.ctx, writeTimeout)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, u, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket at %q: %w", u, err)
	}

	c.conn = conn

	slog.Info("websocket connection established", slog.String("url", u))

	return nil
}

// ReceiveReadyMessage receives the ready message meaning the game is ready to play.
func (c *Client) ReceiveReadyMessage(readyCh chan ReadyMessage) error {
	var msg ReadyMessage

	err := wsjson.Read(c.ctx, c.conn, &msg)
	if err != nil {
		slog.Error("failed to read ready message", slog.Any("error", err))

		return err
	}

	readyCh <- msg

	return nil
}

// ReceiveGameState receives the game state from the server and sends it to the given channel.
func (c *Client) ReceiveGameState(gameStateChan chan<- GameState) {
	defer close(gameStateChan)

	for {
		select {
		case <-c.ctx.Done():
			slog.Info("client context canceled, closing message handler")
		default:
			var gameState GameState
			if err := wsjson.Read(c.ctx, c.conn, &gameState); err != nil {
				slog.Error("failed to read game state: %w", slog.Any("error", err))
			}

			gameStateChan <- gameState
		}
	}
}

// Close closes the WebSocket connection.
func (c *Client) Close() {
	c.cancel()

	if c.conn == nil {
		return
	}

	if err := c.conn.Close(websocket.StatusNormalClosure, "normal closure"); err != nil {
		slog.Error("failed to close websocket connection", slog.Any("error", err))
		return
	}

	slog.Info("websocket connection gracefully closed")
}

// SendPlayerInfo sends the player info to the server.
// It's used to register the player in the server.
func (c *Client) SendPlayerInfo(gi GameInfo) error {
	ctx, cancel := context.WithTimeout(c.ctx, writeTimeout)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, gi); err != nil {
		return fmt.Errorf("failed to send player info: %w", err)
	}

	slog.Info("player info sent successfully")

	return nil
}

// SendPlayerInput sends the player input to the server.
func (c *Client) SendPlayerInput(input PlayerInput) error {
	ctx, cancel := context.WithTimeout(c.ctx, writeTimeout)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, input); err != nil {
		return fmt.Errorf("failed to send player input: %w", err)
	}

	return nil
}
