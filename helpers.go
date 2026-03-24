package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func boolEnabled(v bool) string {
	if v {
		return "enabled"
	}
	return "disabled"
}

func boolYesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func containsAny(s string, parts ...string) bool {
	for _, p := range parts {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

func readRecentLauncherExecLines(logPath string, limit int) ([]string, error) {
	if strings.TrimSpace(logPath) == "" {
		return nil, fmt.Errorf("log path is empty")
	}
	if limit <= 0 {
		return []string{}, nil
	}

	current, err := readRecentMatchingLogLines(logPath, "Launcher exec:", limit)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		current = []string{}
	}

	if len(current) >= limit {
		return current, nil
	}

	remaining := limit - len(current)
	older, oldErr := readRecentMatchingLogLines(logPath+".old", "Launcher exec:", remaining)
	if oldErr != nil {
		if os.IsNotExist(oldErr) {
			return current, nil
		}
		return nil, oldErr
	}

	// Keep chronological order: older log lines first, then current log lines.
	return append(older, current...), nil
}

func readRecentMatchingLogLines(path, needle string, limit int) ([]string, error) {
	if limit <= 0 {
		return []string{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	buf := make([]string, limit)
	count := 0
	next := 0

	for scanner.Scan() {
		line := strings.TrimSpace(strings.TrimSuffix(scanner.Text(), "\r"))
		if line == "" || !strings.Contains(line, needle) {
			continue
		}
		buf[next] = line
		next = (next + 1) % limit
		if count < limit {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if count == 0 {
		return []string{}, nil
	}

	start := 0
	if count == limit {
		start = next
	}

	matches := make([]string, 0, count)
	for i := 0; i < count; i++ {
		idx := (start + i) % limit
		matches = append(matches, buf[idx])
	}
	return matches, nil
}
