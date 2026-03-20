package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/illwill/cardbot/internal/config"
)

// SetupPrompter reads setup answers from a shared buffered input stream.
// Reusing one reader prevents buffered read-ahead from consuming subsequent answers.
type SetupPrompter struct {
	reader *bufio.Reader
	out    io.Writer
}

// NewSetupPrompter creates a setup prompter with shared input/output.
func NewSetupPrompter(in io.Reader, out io.Writer) *SetupPrompter {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	return &SetupPrompter{reader: bufio.NewReader(in), out: out}
}

func (p *SetupPrompter) PromptNamingMode(defaultMode string) string {
	return promptNamingModeReader(p.reader, p.out, defaultMode)
}

func (p *SetupPrompter) PromptDaemonEnabled(defaultEnabled bool) bool {
	return promptDaemonEnabledReader(p.reader, p.out, defaultEnabled)
}

func (p *SetupPrompter) PromptDaemonStartAtLogin(defaultEnabled bool) bool {
	return promptDaemonStartAtLoginReader(p.reader, p.out, defaultEnabled)
}

// RunSetup executes first-time/--setup prompts and persists config.
func RunSetup(
	cfg *config.Config,
	cfgPath string,
	promptDestinationFn func(string) string,
	promptNamingFn func(string) string,
	promptDaemonEnabledFn func(bool) bool,
	promptDaemonStartAtLoginFn func(bool) bool,
) error {
	cfg.Destination.Path = config.ContractPath(promptDestinationFn(cfg.Destination.Path))
	cfg.Naming.Mode = config.NormalizeNamingMode(promptNamingFn(cfg.Naming.Mode))
	cfg.Daemon.Enabled = promptDaemonEnabledFn(cfg.Daemon.Enabled)
	// Keep daemon setup intentionally simple: CardBot now uses the system
	// default terminal app for auto-launch and does not prompt for terminal selection.
	cfg.Daemon.StartAtLogin = promptDaemonStartAtLoginFn(cfg.Daemon.StartAtLogin)
	cfg.Daemon.TerminalApp = "Terminal"
	if !cfg.Daemon.Enabled {
		cfg.Daemon.StartAtLogin = false
	}

	if cfgPath == "" {
		return nil
	}
	return config.Save(cfg, cfgPath)
}

// PromptNamingMode asks the user how filenames should be written on copy.
func PromptNamingMode(defaultMode string) string {
	return promptNamingModeIO(os.Stdin, os.Stdout, defaultMode)
}

// PromptDaemonEnabled asks whether CardBot should auto-launch from daemon mode.
func PromptDaemonEnabled(defaultEnabled bool) bool {
	return promptDaemonEnabledIO(os.Stdin, os.Stdout, defaultEnabled)
}

// PromptDaemonStartAtLogin asks whether daemon mode should auto-start at login.
func PromptDaemonStartAtLogin(defaultEnabled bool) bool {
	return promptDaemonStartAtLoginIO(os.Stdin, os.Stdout, defaultEnabled)
}

func promptNamingModeIO(in io.Reader, out io.Writer, defaultMode string) string {
	return promptNamingModeReader(bufio.NewReader(in), out, defaultMode)
}

func promptNamingModeReader(reader *bufio.Reader, out io.Writer, defaultMode string) string {
	mode := config.NormalizeNamingMode(defaultMode)

	for {
		fmt.Fprintln(out, "────────────────────────────────────────")
		fmt.Fprintln(out, "File Naming")
		fmt.Fprintln(out, "────────────────────────────────────────")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Camera filenames reset every 10,000 shots.")
		fmt.Fprintln(out, "This can cause duplicates when copying multiple cards.")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "[1] Keep camera filenames")
		fmt.Fprintln(out, "    DSC_0001.NEF, DSC_0002.NEF ...")
		fmt.Fprintln(out, "    Use if you rely on camera numbering.")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "[2] Timestamp + sequence")
		fmt.Fprintln(out, "    260314T143052_0001.NEF, _0002.NEF ...")
		fmt.Fprintln(out, "    Use for automatic order across all cards.")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "You can change this later with cardbot --setup.")
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Choice [%s]: ", namingChoiceDefault(mode))

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(out)
			return mode
		}
		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Fprintf(out, "Naming set to: %s\n", namingModeLabel(mode))
			return mode
		}

		if chosen, ok := parseNamingChoice(line); ok {
			fmt.Fprintf(out, "Naming set to: %s\n", namingModeLabel(chosen))
			return chosen
		}

		fmt.Fprintln(out, "Please enter 1 or 2.")
		fmt.Fprintln(out)
	}
}

func promptDaemonEnabledIO(in io.Reader, out io.Writer, defaultEnabled bool) bool {
	return promptDaemonEnabledReader(bufio.NewReader(in), out, defaultEnabled)
}

func promptDaemonEnabledReader(reader *bufio.Reader, out io.Writer, defaultEnabled bool) bool {
	defaultChoice := "n"
	if defaultEnabled {
		defaultChoice = "y"
	}

	for {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Background Auto-Launch")
		fmt.Fprintln(out, "────────────────────────────────────────")
		fmt.Fprintln(out, "When daemon mode is running, launch CardBot")
		fmt.Fprintln(out, "automatically when a memory card is connected? [y/n]")
		fmt.Fprintf(out, "Choice [%s]: ", defaultChoice)

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(out)
			return defaultEnabled
		}
		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Fprintf(out, "Background auto-launch: %s\n", enabledLabel(defaultEnabled))
			return defaultEnabled
		}

		if enabled, ok := parseYesNo(line); ok {
			fmt.Fprintf(out, "Background auto-launch: %s\n", enabledLabel(enabled))
			return enabled
		}

		fmt.Fprintln(out, "Please enter y or n.")
	}
}

func promptDaemonStartAtLoginIO(in io.Reader, out io.Writer, defaultEnabled bool) bool {
	return promptDaemonStartAtLoginReader(bufio.NewReader(in), out, defaultEnabled)
}

func promptDaemonStartAtLoginReader(reader *bufio.Reader, out io.Writer, defaultEnabled bool) bool {
	defaultChoice := "n"
	if defaultEnabled {
		defaultChoice = "y"
	}

	for {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Start at Login")
		fmt.Fprintln(out, "────────────────────────────────────────")
		fmt.Fprintln(out, "Start CardBot daemon automatically")
		fmt.Fprintln(out, "when you log in? [y/n]")
		fmt.Fprintf(out, "Choice [%s]: ", defaultChoice)

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(out)
			return defaultEnabled
		}
		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Fprintf(out, "Start-at-login: %s\n", enabledLabel(defaultEnabled))
			return defaultEnabled
		}

		if enabled, ok := parseYesNo(line); ok {
			fmt.Fprintf(out, "Start-at-login: %s\n", enabledLabel(enabled))
			return enabled
		}

		fmt.Fprintln(out, "Please enter y or n.")
	}
}

func parseYesNo(input string) (bool, bool) {
	switch strings.TrimSpace(strings.ToLower(input)) {
	case "y", "yes":
		return true, true
	case "n", "no":
		return false, true
	default:
		return false, false
	}
}

func enabledLabel(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func parseNamingChoice(input string) (string, bool) {
	switch strings.TrimSpace(strings.ToLower(input)) {
	case "1", "o", "original":
		return config.NamingOriginal, true
	case "2", "t", "timestamp":
		return config.NamingTimestamp, true
	default:
		return "", false
	}
}

func namingChoiceDefault(mode string) string {
	if config.NormalizeNamingMode(mode) == config.NamingTimestamp {
		return "2"
	}
	return "1"
}

func namingModeLabel(mode string) string {
	if config.NormalizeNamingMode(mode) == config.NamingTimestamp {
		return "Timestamp + sequence"
	}
	return "Camera original"
}

func namingDisplayLine(mode string) string {
	if config.NormalizeNamingMode(mode) != config.NamingTimestamp {
		return "Camera original"
	}
	return "Timestamp + sequence"
}
