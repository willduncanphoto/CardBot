//go:build !darwin && !linux

package copy

// diskFreeBytes is not supported on this platform.
func diskFreeBytes(_ string) (int64, bool) { return 0, false }
