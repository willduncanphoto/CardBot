package main

import (
	"bytes"
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
	}

	for _, tt := range tests {
		got, ok := parseNamingChoice(tt.in)
		if ok != tt.wantOK {
			t.Errorf("parseNamingChoice(%q) ok = %v, want %v", tt.in, ok, tt.wantOK)
			continue
		}
		if got != tt.want {
			t.Errorf("parseNamingChoice(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestPromptNamingModeIO_DefaultChoice(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("\n")
	var out bytes.Buffer

	mode := promptNamingModeIO(in, &out, config.NamingOriginal)
	if mode != config.NamingOriginal {
		t.Fatalf("mode = %q, want %q", mode, config.NamingOriginal)
	}
	if !strings.Contains(out.String(), "Choice [1]:") {
		t.Fatalf("expected default choice prompt, got output:\n%s", out.String())
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
		t.Fatalf("expected invalid-input message, got output:\n%s", out.String())
	}
}

func TestNamingStartupLine(t *testing.T) {
	t.Parallel()
	if got := namingStartupLine(config.NamingOriginal); got != "Keep original filenames" {
		t.Fatalf("namingStartupLine(original) = %q", got)
	}
	want := "Timestamp + sequence filenames (YYMMDDTHHMMSS_NNN.EXT, auto 3/4/5 digits)"
	if got := namingStartupLine(config.NamingTimestamp); got != want {
		t.Fatalf("namingStartupLine(timestamp) = %q, want %q", got, want)
	}
}

func TestNamingDisplayLine(t *testing.T) {
	t.Parallel()
	if got := namingDisplayLine(config.NamingOriginal, 3048); got != "Keep original (DSC_0001.NEF)" {
		t.Fatalf("namingDisplayLine(original) = %q", got)
	}
	got := namingDisplayLine(config.NamingTimestamp, 3048)
	want := "Timestamp sequence (260314T143052_0001.NEF) [4-digit]"
	if got != want {
		t.Fatalf("namingDisplayLine(timestamp) = %q, want %q", got, want)
	}
}
