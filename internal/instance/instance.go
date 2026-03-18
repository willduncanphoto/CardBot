package instance

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type pgrepRunner func(name string, args ...string) ([]byte, error)

// HasOtherProcess reports whether any process with the given name exists
// besides selfPID. It uses `pgrep -x <processName>`.
func HasOtherProcess(processName string, selfPID int) (bool, error) {
	return hasOtherProcessWithRunner(processName, selfPID, runPgrep)
}

func hasOtherProcessWithRunner(processName string, selfPID int, run pgrepRunner) (bool, error) {
	processName = strings.TrimSpace(processName)
	if processName == "" {
		return false, fmt.Errorf("process name is required")
	}

	out, err := run("pgrep", "-x", processName)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// pgrep exit 1: no matches.
			return false, nil
		}
		return false, fmt.Errorf("pgrep -x %s failed: %w", processName, err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pid, convErr := strconv.Atoi(line)
		if convErr != nil {
			continue
		}
		if pid != selfPID {
			return true, nil
		}
	}

	return false, nil
}

func runPgrep(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}
