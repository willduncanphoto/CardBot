package launcher

import (
	"fmt"
	"os/exec"
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
		openArgs := append([]string{"-a", app, "--args"}, resolved...)
		return run("open", openArgs...)
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
		return "Terminal"
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

func isGhosttyApp(app string) bool {
	return strings.Contains(strings.ToLower(app), "ghostty")
}

func shQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
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
