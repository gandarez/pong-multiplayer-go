package network

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// NewSpectatorClient creates a new spectator client.
func NewSpectatorClient(ctx context.Context, cancel context.CancelFunc) *Client {
	return &Client{
		ctx:       ctx,
		cancel:    cancel,
		serverURL: BaseURL,
	}
}

// ConnectAsSpectator connects to the server as a spectator.
func (c *Client) ConnectAsSpectator(sessionID string) error {
	u := fmt.Sprintf("wss://%s/spectate", c.serverURL)

	ctx, cancel := context.WithTimeout(c.ctx, writeTimeout)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, u, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket at %q: %w", u, err)
	}

	c.conn = conn

	// send spectate request with session ID
	spectateRequest := map[string]string{
		"session_id": sessionID,
	}
	if err := wsjson.Write(ctx, c.conn, spectateRequest); err != nil {
		return fmt.Errorf("failed to send spectate request: %w", err)
	}

	slog.Info("websocket connection established as spectator", slog.String("url", u))

	return nil
}
