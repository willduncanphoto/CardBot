//go:build darwin

package pick

import (
	"fmt"
	"os/exec"
	"strings"
)

// Folder opens the native macOS folder picker dialog and returns the chosen path.
func Folder(defaultPath string) (string, error) {
	script := fmt.Sprintf(
		`POSIX path of (choose folder with prompt "Where should CardBot copy your work?" default location POSIX file "%s")`,
		defaultPath,
	)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
