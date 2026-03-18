package main

import (
	"errors"
	"strings"
	"testing"
)

func TestDaemonLaunchHint_Automation(t *testing.T) {
	hint := daemonLaunchHint(errors.New("not authorized to send Apple events to Terminal"))
	if !strings.Contains(strings.ToLower(hint), "automation") {
		t.Fatalf("hint = %q, expected automation guidance", hint)
	}
}

func TestDaemonLaunchHint_FullDiskAccess(t *testing.T) {
	hint := daemonLaunchHint(errors.New("operation not permitted"))
	if !strings.Contains(strings.ToLower(hint), "full disk access") {
		t.Fatalf("hint = %q, expected full disk access guidance", hint)
	}
}

func TestDaemonLaunchHint_Unknown(t *testing.T) {
	hint := daemonLaunchHint(errors.New("random error"))
	if hint != "" {
		t.Fatalf("hint = %q, want empty", hint)
	}
}
