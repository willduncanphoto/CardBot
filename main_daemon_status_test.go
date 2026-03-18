package main

import "testing"

func TestParseDaemonStatusOptions_Default(t *testing.T) {
	t.Parallel()

	opts, err := parseDaemonStatusOptions(nil)
	if err != nil {
		t.Fatalf("parseDaemonStatusOptions error: %v", err)
	}
	if opts.JSON {
		t.Fatal("opts.JSON = true, want false")
	}
}

func TestParseDaemonStatusOptions_JSON(t *testing.T) {
	t.Parallel()

	opts, err := parseDaemonStatusOptions([]string{"--json"})
	if err != nil {
		t.Fatalf("parseDaemonStatusOptions error: %v", err)
	}
	if !opts.JSON {
		t.Fatal("opts.JSON = false, want true")
	}
}

func TestParseDaemonStatusOptions_UnexpectedArg(t *testing.T) {
	t.Parallel()

	_, err := parseDaemonStatusOptions([]string{"wat"})
	if err == nil {
		t.Fatal("expected error")
	}
}
