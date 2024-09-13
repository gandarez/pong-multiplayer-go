package level

// Level represents the level of the game.
type Level int

const (
	// Easy represents the easy level.
	Easy Level = iota
	// Medium represents the medium level.
	Medium
	// Hard represents the hard level.
	Hard
)

// String returns a string representation of the level.
func (l Level) String() string {
	return [...]string{"Easy", "Medium", "Hard"}[l]
}
