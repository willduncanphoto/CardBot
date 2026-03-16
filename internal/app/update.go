package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/illwill/cardbot/internal/config"
	cblog "github.com/illwill/cardbot/internal/log"
	"github.com/illwill/cardbot/internal/update"
)

const (
	updateCheckInterval = 24 * time.Hour // Check for updates every 24 hours
	updateCheckTimeout  = 5 * time.Second
	selfUpdateTimeout   = 60 * time.Second
)

// MaybeCheckForUpdate checks for updates if the cache interval has elapsed.
func MaybeCheckForUpdate(cfg *config.Config, cfgPath string, logger *cblog.Logger, version string) (string, bool) {
	// Check if enough time has elapsed since last check
	if !update.ShouldCheck(cfg.Update.LastCheck, time.Now(), updateCheckInterval) {
		if logger != nil {
			logger.Printf("Update check skipped (cached)")
		}
		return "", false
	}

	fmt.Printf("[%s] Checking for updates\n", ts())

	ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
	defer cancel()
	res, err := update.CheckLatest(ctx, nil, update.DefaultAPIBase, update.DefaultRepo, version)
	if err != nil {
		if logger != nil {
			logger.Printf("Update check failed: %v", err)
		}
		return "", false
	}

	// Update the last check timestamp in config and save it
	cfg.Update.LastCheck = time.Now().Format(time.RFC3339)
	if cfgPath != "" {
		if saveErr := config.Save(cfg, cfgPath); saveErr != nil && logger != nil {
			logger.Printf("Warning: could not save update check timestamp: %v", saveErr)
		}
	}

	if res.Update {
		return res.Latest, true
	}
	return "", false
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
