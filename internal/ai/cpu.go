package ai

import "math/rand/v2"

const cpuSpeed = 4

// GuessBallPosition returns the new position of the enemy paddle based on the ball position.
// It returns the new Y position of the enemy paddle.
func GuessBallPosition(ballY, enemyY, enemyHeight, screenHeight, fieldBorderWidth float64) float64 {
	delta := float64(rand.IntN(15)) // nolint:gosec

	if enemyY < ballY-delta {
		enemyY += cpuSpeed // Move down
	}

	if enemyY > ballY+delta {
		enemyY -= cpuSpeed // Move up
	}

	return keepInBounds(enemyY, enemyHeight, screenHeight, fieldBorderWidth)
}

func keepInBounds(y, height, screenHeight, fieldBorderWidth float64) float64 {
	if y < 0 {
		return 0
	}

	if y > screenHeight-height-fieldBorderWidth {
		return screenHeight - height - fieldBorderWidth
	}

	return y
}
