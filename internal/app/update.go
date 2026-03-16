package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/illwill/cardbot/internal/config"
	cblog "github.com/illwill/cardbot/internal/log"
	"github.com/illwill/cardbot/internal/update"
)

const (
	updateCheckTimeout = 5 * time.Second
	selfUpdateTimeout  = 60 * time.Second
)

var httpCodeRe = regexp.MustCompile(`http (\d{3})`)

// MaybeCheckForUpdate checks for updates on every app startup.
func MaybeCheckForUpdate(cfg *config.Config, cfgPath string, logger *cblog.Logger, version string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
	defer cancel()
	res, err := update.CheckLatest(ctx, nil, update.DefaultAPIBase, update.DefaultRepo, version)
	if err != nil {
		if logger != nil {
			logger.Printf("Update check failed: %v", err)
		}
		return "", err
	}

	if res.Update {
		return res.Latest, nil
	}
	return "", nil
}

// UpdateErrCode returns a short display code for an update check failure,
// or an empty string if the error is a plain connectivity loss (no code needed).
func UpdateErrCode(err error) string {
	if err == nil {
		return ""
	}
	// Deadline elapsed — plain no-signal, no code.
	if errors.Is(err, context.DeadlineExceeded) {
		return ""
	}
	s := strings.ToLower(err.Error())
	// Network-level failures — no code needed.
	if strings.Contains(s, "no such host") ||
		strings.Contains(s, "no route to host") ||
		strings.Contains(s, "network is unreachable") ||
		strings.Contains(s, "i/o timeout") ||
		strings.Contains(s, "connection refused") {
		return ""
	}
	// HTTP status code — extract and return just the number.
	if m := httpCodeRe.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return "ERR"
}

// RunSelfUpdate performs a self-update to the latest version.
func RunSelfUpdate(version string) int {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not determine executable path: %s\n", friendlyErr(err))
		return 1
	}

	fmt.Printf("[%s] Checking for updates…\n", ts())
	ctx, cancel := context.WithTimeout(context.Background(), selfUpdateTimeout)
	defer cancel()

	installed, err := update.SelfUpdate(ctx, nil, update.DefaultAPIBase, update.DefaultRepo, version, execPath)
	if err == nil {
		fmt.Printf("[%s] Updated successfully to %s\n", ts(), installed)
		fmt.Printf("[%s] Restart CardBot to use the new version.\n", ts())
		return 0
	}

	if errors.Is(err, update.ErrAlreadyUpToDate) {
		fmt.Printf("[%s] CardBot is already up to date (%s)\n", ts(), version)
		return 0
	}

	fmt.Fprintf(os.Stderr, "Error: %s\n", friendlyErr(err))
	if isPermissionErr(err) {
		fmt.Fprintf(os.Stderr, "Try: sudo %q self-update\n", execPath)
	}
	return 1
}

func isPermissionErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, os.ErrPermission) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "permission denied")
}
