package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Options controls how a terminal is launched for a detected card.
type Options struct {
	TerminalApp   string
	LaunchArgs    []string
	CardBotBinary string
	MountPath     string
}

type commandRunner func(name string, args ...string) error

// Launch opens the configured terminal and runs cardbot for the given mount path.
func Launch(opts Options) error {
	return launchWith(opts, runCommand)
}

func launchWith(opts Options, run commandRunner) error {
	binary := strings.TrimSpace(opts.CardBotBinary)
	mountPath := strings.TrimSpace(opts.MountPath)
	if binary == "" {
		return fmt.Errorf("cardbot binary path is required")
	}
	if mountPath == "" {
		return fmt.Errorf("mount path is required")
	}

	app := normalizeTerminalApp(opts.TerminalApp)

	if len(opts.LaunchArgs) > 0 {
		resolved := resolveLaunchArgs(opts.LaunchArgs, binary, mountPath)
		if isSystemDefaultTerminal(app) {
			return run("open", resolved...)
		}
		openArgs := append([]string{"-a", app, "--args"}, resolved...)
		return run("open", openArgs...)
	}

	if isSystemDefaultTerminal(app) {
		scriptPath, err := writeDefaultTerminalCommandScript(binary, mountPath)
		if err != nil {
			return err
		}
		return run("open", scriptPath)
	}

	if isTerminalApp(app) {
		cmd := fmt.Sprintf("%s %s", shQuote(binary), shQuote(mountPath))
		return run("osascript",
			"-e", fmt.Sprintf(`tell application "Terminal" to do script %q`, cmd),
			"-e", `activate application "Terminal"`,
		)
	}

	if isGhosttyApp(app) {
		cmd := fmt.Sprintf("%s %s", shQuote(binary), shQuote(mountPath))
		return run("open", "-a", app, "--args", "-e", cmd)
	}

	return run("open", "-a", app, "--args", binary, mountPath)
}

func normalizeTerminalApp(app string) string {
	app = strings.TrimSpace(app)
	if app == "" {
		return "Default"
	}
	if strings.EqualFold(app, "terminal.app") {
		return "Terminal"
	}
	if strings.EqualFold(app, "default") || strings.EqualFold(app, "system default") || strings.EqualFold(app, "macos default") {
		return "Default"
	}
	if strings.EqualFold(app, "ghostty") {
		return "Ghostty"
	}
	return app
}

func resolveLaunchArgs(args []string, binary, mountPath string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		replaced := strings.ReplaceAll(arg, "{{mount_path}}", mountPath)
		replaced = strings.ReplaceAll(replaced, "{{cardbot_binary}}", binary)
		out = append(out, replaced)
	}
	return out
}

func isTerminalApp(app string) bool {
	a := strings.ToLower(strings.TrimSpace(app))
	return a == "terminal" || a == "terminal.app"
}

func isSystemDefaultTerminal(app string) bool {
	a := strings.ToLower(strings.TrimSpace(app))
	return a == "default" || a == "system default" || a == "macos default"
}

func isGhosttyApp(app string) bool {
	return strings.Contains(strings.ToLower(app), "ghostty")
}

func shQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func writeDefaultTerminalCommandScript(binary, mountPath string) (string, error) {
	f, err := os.CreateTemp("", "cardbot-launch-*.command")
	if err != nil {
		return "", fmt.Errorf("creating command script: %w", err)
	}
	defer f.Close()

	scriptPath := f.Name()
	script := fmt.Sprintf("#!/bin/sh\nrm -- %s\nexec %s %s\n", shQuote(scriptPath), shQuote(binary), shQuote(mountPath))
	if _, err := f.WriteString(script); err != nil {
		_ = os.Remove(scriptPath)
		return "", fmt.Errorf("writing command script: %w", err)
	}
	if err := f.Chmod(0o700); err != nil {
		_ = os.Remove(scriptPath)
		return "", fmt.Errorf("chmod command script: %w", err)
	}
	return filepath.Clean(scriptPath), nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
		}
		return fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, msg)
	}
	return nil
}
