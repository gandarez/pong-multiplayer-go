package level

type Level int

const (
	Easy Level = iota
	Medium
	Hard
)

func (l Level) String() string {
	return [...]string{"Easy", "Medium", "Hard"}[l]
}
