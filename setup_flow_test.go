package main

import (
	"path/filepath"
	"testing"

	"github.com/illwill/cardbot/internal/config"
)

func TestRunSetup_WritesNamingModeToConfig(t *testing.T) {
	t.Parallel()

	cfg := config.Defaults()
	cfgPath := filepath.Join(t.TempDir(), "config.json")

	err := runSetup(
		cfg,
		cfgPath,
		func(string) string { return "~/Pictures/Jobs" },
		func(string) string { return config.NamingTimestamp },
	)
	if err != nil {
		t.Fatalf("runSetup error: %v", err)
	}

	loaded, warnings, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if loaded.Destination.Path != "~/Pictures/Jobs" {
		t.Fatalf("Destination.Path = %q, want %q", loaded.Destination.Path, "~/Pictures/Jobs")
	}
	if loaded.Naming.Mode != config.NamingTimestamp {
		t.Fatalf("Naming.Mode = %q, want %q", loaded.Naming.Mode, config.NamingTimestamp)
	}
}

func TestRunSetup_NormalizesInvalidNamingMode(t *testing.T) {
	t.Parallel()

	cfg := config.Defaults()
	cfgPath := filepath.Join(t.TempDir(), "config.json")

	err := runSetup(
		cfg,
		cfgPath,
		func(string) string { return "~/Pictures/Jobs" },
		func(string) string { return "banana" },
	)
	if err != nil {
		t.Fatalf("runSetup error: %v", err)
	}

	loaded, _, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if loaded.Naming.Mode != config.NamingOriginal {
		t.Fatalf("Naming.Mode = %q, want %q", loaded.Naming.Mode, config.NamingOriginal)
	}
}
