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

type Client struct {
	conn      *websocket.Conn
	serverURL string
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClient(ctx context.Context, cancel context.CancelFunc, serverURL string) *Client {
	return &Client{
		serverURL: serverURL,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (c *Client) Connect() error {
	u := fmt.Sprintf("wss://%s/multiplayer", c.serverURL)

	ctx, cancel := context.WithTimeout(c.ctx, writeTimeout)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, u, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket at %s: %w", u, err)
	}

	c.conn = conn

	slog.Info("websocket connection established", slog.String("url", u))

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

func (c *Client) Close() {
	c.cancel()

	if c.conn == nil {
		return
	}

	if err := c.conn.Close(websocket.StatusNormalClosure, "normal closure"); err != nil {
		slog.Error("error closing WebSocket connection", slog.Any("error", err))
	} else {
		slog.Info("websocket connection closed gracefully")
	}
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

	slog.Info("player input sent successfully")

	return nil
}
