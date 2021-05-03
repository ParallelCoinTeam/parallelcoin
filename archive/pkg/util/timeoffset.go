package util

import "time"

// GetSecondsAhead returns the difference in time, of the second ahead of the first
func GetSecondsAhead(first, second time.Time) int {
	return int(second.Sub(first).Truncate(time.Second) / time.Second)
}
