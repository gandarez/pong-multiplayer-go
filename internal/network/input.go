package network

// PlayerInput represents the input of the player when it is sent over the network.
type PlayerInput struct {
	Up   bool `json:"up"`
	Down bool `json:"down"`
}
