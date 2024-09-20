package network

import "github.com/gandarez/pong-multiplayer-go/pkg/geometry"

// GameState represents the state of the game when it is sent over the network.
type GameState struct {
	BallPosition   geometry.Vector `json:"ball_position"`
	CurrentPlayer  PlayerState     `json:"current"`
	OpponentPlayer PlayerState     `json:"opponent"`
}

// PlayerState represents the state of a player when it is sent over the network.
type PlayerState struct {
	// TODO: add player ID to better identify them
	PositionY float64       `json:"position_y"`
	Side      geometry.Side `json:"side"`
	Score     int8          `json:"score"`
	Ping      int64         `json:"ping"`
}
