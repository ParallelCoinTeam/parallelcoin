package limits

// SetLimits is a no-op on Plan 9 due to the lack of process accounting.
func SetLimits() error {
	return nil
}
