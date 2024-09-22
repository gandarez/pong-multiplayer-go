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
	BaseURL     = "game.go-go.dev"
	writeWait   = 10 * time.Second
	readTimeout = 60 * time.Second
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

	ctx, cancel := context.WithTimeout(c.ctx, writeWait)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, u, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket at %s: %w", u, err)
	}

	c.conn = conn
	slog.Info("WebSocket connection established", slog.String("url", u))
	return nil
}

func (c *Client) ReceiveGameState(gameStateChan chan<- GameState) error {
	defer close(gameStateChan)

	for {
		select {
		case <-c.ctx.Done():
			slog.Info("Client context canceled, closing message handler")
			return nil
		default:
			ctx, cancel := context.WithTimeout(c.ctx, readTimeout)
			defer cancel()

			var gameState GameState
			if err := wsjson.Read(ctx, c.conn, &gameState); err != nil {
				return fmt.Errorf("failed to read game state: %w", err)
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
		slog.Error("Error closing WebSocket connection", slog.Any("error", err))
	} else {
		slog.Info("WebSocket connection closed gracefully")
	}
}

func (c *Client) SendPlayerInfo(pi PlayerInfo) error {
	ctx, cancel := context.WithTimeout(c.ctx, writeWait)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, pi); err != nil {
		return fmt.Errorf("failed to send player info: %w", err)
	}

	slog.Info("Player info sent successfully")
	return nil
}

func (c *Client) SendPlayerInput(input PlayerInput) error {
	ctx, cancel := context.WithTimeout(c.ctx, writeWait)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, input); err != nil {
		return fmt.Errorf("failed to send player input: %w", err)
	}

	slog.Info("Player input sent successfully")
	return nil
}
