package copy

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	namingModeOriginal  = "original"
	namingModeTimestamp = "timestamp"
)

func normalizeNamingMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case namingModeTimestamp:
		return namingModeTimestamp
	case namingModeOriginal, "":
		return namingModeOriginal
	default:
		return namingModeOriginal
	}
}

// sequenceDigits picks minimum padding for card file count, clamped to 3..5.
func sequenceDigits(totalFiles int) int {
	switch {
	case totalFiles <= 999:
		return 3
	case totalFiles <= 9999:
		return 4
	default:
		return 5
	}
}

func sequenceMax(digits int) int {
	switch digits {
	case 3:
		return 999
	case 4:
		return 9999
	default:
		return 99999
	}
}

func formatSequence(n, digits int) string {
	if n < 1 {
		n = 1
	}
	if digits < 3 {
		digits = 3
	}
	if digits > 5 {
		digits = 5
	}
	return fmt.Sprintf("%0*d", digits, n)
}

func timestampStem(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}
	return t.Format("060102T150405")
}

func renamedRelativePath(relPath string, captureTime time.Time, seq, digits int) string {
	dir := filepath.Dir(relPath)
	ext := strings.ToUpper(filepath.Ext(relPath))
	name := timestampStem(captureTime) + "_" + formatSequence(seq, digits) + ext
	if dir == "." || dir == "" {
		return name
	}
	return filepath.Join(dir, name)
}
