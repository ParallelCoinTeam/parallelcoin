package limits

// SetLimits is a no-op on Windows since it's not required there.
func SetLimits() (e error) {
	return nil
}
