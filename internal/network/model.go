package network

import "github.com/gandarez/pong-multiplayer-go/pkg/geometry"

type (
	// GameState represents the state of the game when it is sent over the network.
	GameState struct {
		Ball           BallState   `json:"ball"`
		CurrentPlayer  PlayerState `json:"current"`
		OpponentPlayer PlayerState `json:"opponent"`
	}

	// BallState represents the state of the ball when it is sent over the network.
	BallState struct {
		Angle    float64         `json:"angle"`
		Bounces  int             `json:"bounces"`
		Position geometry.Vector `json:"position"`
	}

	// PlayerState represents the state of a player when it is sent over the network.
	PlayerState struct {
		Name      string        `json:"name"`
		PositionY float64       `json:"position_y"`
		Side      geometry.Side `json:"side"`
		Score     int8          `json:"score"`
		Ping      int64         `json:"ping"`
	}

	// GameInfo represents the information of a game when it is sent over the network.
	GameInfo struct {
		PlayerName       string `json:"player_name"`
		Level            int    `json:"level"`
		ScreenWidth      int    `json:"screen_width"`
		ScreenHeight     int    `json:"screen_height"`
		MaxScore         int    `json:"max_score"`
		FieldBorderWidth int    `json:"field_border_width"`
	}

	// PlayerInput represents the keyboard/touch input of the player when it is sent over the network.
	PlayerInput struct {
		Up   bool `json:"up"`
		Down bool `json:"down"`
	}
)
