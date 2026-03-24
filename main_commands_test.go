package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunNoArgSubcommand_RejectsExtraArgs(t *testing.T) {
	t.Parallel()

	runCalled := false
	code := runNoArgSubcommand("install-daemon", []string{"extra"}, func() int {
		runCalled = true
		return 0
	})
	if code != 2 {
		t.Fatalf("runNoArgSubcommand() code = %d, want 2", code)
	}
	if runCalled {
		t.Fatal("runNoArgSubcommand() executed callback with extra args")
	}
}

func TestLooksLikeCommandToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "unknown command token", arg: "daemon-statuz", want: true},
		{name: "flag", arg: "--setup", want: false},
		{name: "absolute path", arg: "/Volumes/NIKON", want: false},
		{name: "relative path", arg: "./card", want: false},
		{name: "home path", arg: "~/card", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := looksLikeCommandToken(tt.arg); got != tt.want {
				t.Fatalf("looksLikeCommandToken(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestLooksLikeCommandToken_ExistingRelativePath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	temp := t.TempDir()
	if err := os.Chdir(temp); err != nil {
		t.Fatalf("Chdir temp: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	name := "card"
	if err := os.Mkdir(filepath.Join(temp, name), 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	if got := looksLikeCommandToken(name); got {
		t.Fatalf("looksLikeCommandToken(%q) = true, want false for existing path", name)
	}
}

func TestTryRunSubcommand_UnknownCommand(t *testing.T) {
	t.Parallel()

	handled, code := tryRunSubcommand([]string{"daemon-statuz"})
	if !handled {
		t.Fatal("handled = false, want true")
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestTryRunSubcommand_PathArgumentFallsThrough(t *testing.T) {
	t.Parallel()

	handled, code := tryRunSubcommand([]string{"/Volumes/NIKON"})
	if handled {
		t.Fatal("handled = true, want false")
	}
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
}
