package sanitizers

import (
	"time"
)

func ClampInt(min, max int) func(input int) (result int) {
	return func(input int) (result int) {
		if input > max {
			return max
		}
		if input < min {
			return min
		}
		return input
	}
}

func ClampFloat(min, max float64) func(input float64) (result float64) {
	return func(input float64) (result float64) {
		if input > max {
			return max
		}
		if input < min {
			return min
		}
		return input
	}
}

func ClampDuration(min, max time.Duration) func(input time.Duration) (result time.Duration) {
	return func(input time.Duration) (result time.Duration) {
		if input > max {
			return max
		}
		if input < min {
			return min
		}
		return input
	}
}
