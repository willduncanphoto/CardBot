package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/illwill/cardbot/internal/config"
)

// promptNamingMode asks the user how filenames should be written on copy.
func promptNamingMode(defaultMode string) string {
	return promptNamingModeIO(os.Stdin, os.Stdout, defaultMode)
}

func promptNamingModeIO(in io.Reader, out io.Writer, defaultMode string) string {
	mode := config.NormalizeNamingMode(defaultMode)
	reader := bufio.NewReader(in)

	for {
		fmt.Fprintln(out, "────────────────────────────────────────")
		fmt.Fprintln(out, "File Naming")
		fmt.Fprintln(out, "────────────────────────────────────────")
		fmt.Fprintln(out, "How would you like files named when copying?")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "[1] Keep original filenames (DSC_0001.NEF)")
		fmt.Fprintln(out, "[2] Timestamp + sequence (260314T143052_001.NEF)")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "The timestamp comes from when each photo was taken.")
		fmt.Fprintln(out, "Sequence digits adjust automatically based on card size (3/4/5 digits).")
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
		return "Timestamp + sequence (auto-detected digits)"
	}
	return "Keep original filenames"
}

func namingStartupLine(mode string) string {
	if config.NormalizeNamingMode(mode) == config.NamingTimestamp {
		return "Timestamp + sequence filenames (YYMMDDTHHMMSS_NNN.EXT, auto 3/4/5 digits)"
	}
	return "Keep original filenames"
}

func namingDisplayLine(mode string, fileCount int) string {
	if config.NormalizeNamingMode(mode) != config.NamingTimestamp {
		return "Keep original (DSC_0001.NEF)"
	}
	digits := sequenceDigitsForCount(fileCount)
	sample := "001"
	switch digits {
	case 4:
		sample = "0001"
	case 5:
		sample = "00001"
	}
	return fmt.Sprintf("Timestamp sequence (260314T143052_%s.NEF) [%d-digit]", sample, digits)
}

func sequenceDigitsForCount(fileCount int) int {
	switch {
	case fileCount <= 999:
		return 3
	case fileCount <= 9999:
		return 4
	default:
		return 5
	}
}
