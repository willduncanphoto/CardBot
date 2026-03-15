package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/illwill/cardbot/internal/config"
)

func TestParseNamingChoice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{"1", config.NamingOriginal, true},
		{"original", config.NamingOriginal, true},
		{"o", config.NamingOriginal, true},
		{"2", config.NamingTimestamp, true},
		{"timestamp", config.NamingTimestamp, true},
		{"t", config.NamingTimestamp, true},
		{"nope", "", false},
		{"", "", false},
		{"  2  ", config.NamingTimestamp, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input_%q", tt.in), func(t *testing.T) {
			got, ok := parseNamingChoice(tt.in)
			if ok != tt.wantOK {
				t.Errorf("parseNamingChoice(%q) ok = %v, want %v", tt.in, ok, tt.wantOK)
				return
			}
			if got != tt.want {
				t.Errorf("parseNamingChoice(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPromptNamingModeIO_DefaultOriginal(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("\n")
	var out bytes.Buffer

	mode := promptNamingModeIO(in, &out, config.NamingOriginal)
	if mode != config.NamingOriginal {
		t.Fatalf("mode = %q, want %q", mode, config.NamingOriginal)
	}
	if !strings.Contains(out.String(), "Choice [1]:") {
		t.Fatalf("expected default [1] prompt, got:\n%s", out.String())
	}
}

func TestPromptNamingModeIO_DefaultTimestamp(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("\n")
	var out bytes.Buffer

	mode := promptNamingModeIO(in, &out, config.NamingTimestamp)
	if mode != config.NamingTimestamp {
		t.Fatalf("mode = %q, want %q", mode, config.NamingTimestamp)
	}
	if !strings.Contains(out.String(), "Choice [2]:") {
		t.Fatalf("expected default [2] prompt, got:\n%s", out.String())
	}
}

func TestPromptNamingModeIO_InvalidThenValid(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("x\n2\n")
	var out bytes.Buffer

	mode := promptNamingModeIO(in, &out, config.NamingOriginal)
	if mode != config.NamingTimestamp {
		t.Fatalf("mode = %q, want %q", mode, config.NamingTimestamp)
	}
	if !strings.Contains(out.String(), "Please enter 1 or 2.") {
		t.Fatalf("expected invalid-input message, got:\n%s", out.String())
	}
}

func TestPromptNamingModeIO_EOF(t *testing.T) {
	t.Parallel()
	// Simulate EOF (no input at all).
	in := strings.NewReader("")
	var out bytes.Buffer

	mode := promptNamingModeIO(in, &out, config.NamingTimestamp)
	if mode != config.NamingTimestamp {
		t.Fatalf("mode = %q, want %q (should return default on EOF)", mode, config.NamingTimestamp)
	}
}

func TestNamingStartupLine(t *testing.T) {
	t.Parallel()

	// Note: namingStartupLine is no longer used at startup (0.4.0 UX cleanup)
	// These tests verify the simplified format for potential future use
	t.Run("original", func(t *testing.T) {
		got := namingStartupLine(config.NamingOriginal)
		want := "Camera original (DSC_xxxx.NEF)"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("timestamp", func(t *testing.T) {
		want := "Timestamp + sequence (0001-9999)"
		got := namingStartupLine(config.NamingTimestamp)
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func TestNamingDisplayLine(t *testing.T) {
	t.Parallel()

	t.Run("original", func(t *testing.T) {
		got := namingDisplayLine(config.NamingOriginal, 3048)
		want := "Camera original (DSC_xxxx.NEF)"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("timestamp", func(t *testing.T) {
		got := namingDisplayLine(config.NamingTimestamp, 3048)
		want := "Timestamp + sequence (xxxx = 0001-9999)"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
