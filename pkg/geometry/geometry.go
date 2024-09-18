package geometry

import "fmt"

type (
	// Rect represents a rectangle.
	Rect struct {
		X      float64
		Y      float64
		Width  float64
		Height float64
	}

	// Side represents the side of the board.
	Side int

	// Vector represents a 2D vector.
	Vector struct {
		X float64
		Y float64
	}
)

const (
	// Undefined represents an undefined side.
	Undefined Side = iota
	// Right is the right side of the board.
	Right
	// Left is the left side of the board.
	Left
)

// MaxX returns the maximum X value of the rectangle.
func (r Rect) MaxX() float64 {
	return r.X + r.Width
}

// MaxY returns the maximum Y value of the rectangle.
func (r Rect) MaxY() float64 {
	return r.Y + r.Height
}

// Intersects returns true if the rectangle intersects with another rectangle.
func (r Rect) Intersects(other Rect) bool {
	return r.X <= other.MaxX() &&
		other.X <= r.MaxX() &&
		r.Y <= other.MaxY() &&
		other.Y <= r.MaxY()
}

// String returns a string representation of the rectangle.
func (r Rect) String() string {
	return fmt.Sprintf("x:%.f-y:%.f - w:%.f-h:%.f", r.X, r.Y, r.Width, r.Height)
}

// String returns a string representation of the vector.
func (v Vector) String() string {
	return fmt.Sprintf("%.f:%.f", v.X, v.Y)
}
