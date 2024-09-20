package network

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// BaseURL is the base URL of the server.
const BaseURL = "game.go-go.dev"

// Client represents a client that connects to the server using websocket.
type Client struct {
	conn       *websocket.Conn
	serverURL  string
	playerName string
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewClient creates a new client with the given player name and server URL.
func NewClient(ctx context.Context, cancel context.CancelFunc, playerName, serverURL string) *Client {
	return &Client{
		serverURL:  serverURL,
		playerName: playerName,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Connect connects the client to the server.
func (c *Client) Connect() error {
	u := url.URL{Scheme: "ws", Host: c.serverURL, Path: "/multiplayer"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn

	playerInfo := struct {
		Name string `json:"name"`
	}{
		Name: c.playerName,
	}

	if err := c.conn.WriteJSON(playerInfo); err != nil {
		return err
	}

	c.conn.SetPingHandler(func(appData string) error {
		return c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	return nil
}

// Close closes the client connection.
func (c *Client) Close() {
	c.cancel()

	if c.conn == nil {
		return
	}

	if err := c.conn.Close(); err != nil {
		slog.Error("error closing connection", slog.Any("error", err))
	}
}

// ReceiveGameState receives the game state from the server and sends it to the given channel.
func (c *Client) ReceiveGameState(gameStateChan chan<- GameState) {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				close(gameStateChan)
				return
			default:
				var gameState GameState
				if err := c.conn.ReadJSON(&gameState); err != nil {
					slog.Error("error reading json from connection", slog.Any("error", err))
					continue
				}
				gameStateChan <- gameState
			}
		}
	}()
}

// SendPlayerInput sends the player input to the server.
func (c *Client) SendPlayerInput(input PlayerInput) error {
	return c.conn.WriteJSON(input)
}
