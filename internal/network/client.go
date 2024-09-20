package network

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const BaseURL = "game.go-go.dev"

type Client struct {
	conn       *websocket.Conn
	serverURL  string
	playerName string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewClient(playerName, serverURL string) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		serverURL:  serverURL,
		playerName: playerName,
		ctx:        ctx,
		cancel:     cancel,
	}
}

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

func (c *Client) Close() {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
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
				if err := c.conn.ReadJSON(&gameState); err != nil {
					slog.Error("error reading json from connection", slog.Any("error", err))
					continue
				}
				gameStateChan <- gameState
			}
		}
	}()
}

func (c *Client) SendPlayerInput(input PlayerInput) error {
	return c.conn.WriteJSON(input)
}
