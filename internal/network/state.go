package network

import "github.com/gandarez/pong-multiplayer-go/pkg/geometry"

type GameState struct {
	BallPosition   geometry.Vector `json:"ball_position"`
	CurrentPlayer  PlayerState     `json:"current"`
	OpponentPlayer PlayerState     `json:"opponent"`
}

type PlayerState struct {
	// TODO: add player ID to better identify them
	PositionY float64       `json:"position_y"`
	Side      geometry.Side `json:"side"`

	// TODO: Change to int8 instead of int
	Score int8  `json:"score"`
	Ping  int64 `json:"ping"`
}
