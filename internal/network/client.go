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
	BaseURL     = "localhost:8080"
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
	u := fmt.Sprintf("ws://%s/multiplayer", c.serverURL)

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

// ReceiveReadyMessage listens for the ReadyMessage from the server
func (c *Client) ReceiveReadyMessage(readyCh chan ReadyMessage) error {
	var msg ReadyMessage
	err := wsjson.Read(c.ctx, c.conn, &msg)
	if err != nil {
		slog.Error("Failed to read ReadyMessage", slog.Any("error", err))
		return err
	}

	readyCh <- msg

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
			var gameState GameState
			if err := wsjson.Read(c.ctx, c.conn, &gameState); err != nil {
				slog.Error("Failed to read game state", slog.Any("error", err))
				return fmt.Errorf("failed to read game state: %w", err)
			}

			slog.Info("Received game state", slog.Any("gameState", gameState))

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

func (c *Client) SendPlayerInfo(playerInfo GameInfo) error {
	ctx, cancel := context.WithTimeout(c.ctx, writeWait)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, playerInfo); err != nil {
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
