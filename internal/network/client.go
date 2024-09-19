package network

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn       *websocket.Conn
	serverURL  string
	PlayerName string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewClient(playerName, serverURL string) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		serverURL:  serverURL,
		PlayerName: playerName,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (c *Client) Connect() error {
	u := url.URL{Scheme: "ws", Host: c.serverURL, Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.Conn = conn

	playerInfo := struct {
		Name string `json:"name"`
	}{
		Name: c.PlayerName,
	}

	if err := c.Conn.WriteJSON(playerInfo); err != nil {
		return err
	}

	c.Conn.SetPingHandler(func(appData string) error {
		return c.Conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	return nil
}

func (c *Client) Close() {
	c.cancel()
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *Client) ReceiveGameState(gameStateChan chan<- GameState) {
	go func() {
		defer close(gameStateChan)

		for {
			select {
			case <-c.ctx.Done():
				close(gameStateChan)
				return
			default:
				var gameState GameState
				if err := c.Conn.ReadJSON(&gameState); err != nil {
					slog.Error("Error reading JSON from connection", slog.Any("error", err))
					continue
				}
				gameStateChan <- gameState
			}
		}
	}()
}

func (c *Client) SendPlayerInput(input PlayerInput) error {
	return c.Conn.WriteJSON(input)
}
