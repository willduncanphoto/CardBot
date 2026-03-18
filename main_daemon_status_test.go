package main

import (
	"errors"
	"testing"
)

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

func TestCollectSingleInstanceGuardStatus_OtherProcess(t *testing.T) {
	t.Parallel()

	st := collectSingleInstanceGuardStatus("cardbot", 1234, func(processName string, selfPID int) (bool, error) {
		if processName != "cardbot" {
			t.Fatalf("processName = %q, want %q", processName, "cardbot")
		}
		if selfPID != 1234 {
			t.Fatalf("selfPID = %d, want 1234", selfPID)
		}
		return true, nil
	})

	if !st.Enabled {
		t.Fatal("Enabled = false, want true")
	}
	if !st.HasOtherProcess {
		t.Fatal("HasOtherProcess = false, want true")
	}
	if st.CheckError != "" {
		t.Fatalf("CheckError = %q, want empty", st.CheckError)
	}
}

func TestCollectSingleInstanceGuardStatus_CheckError(t *testing.T) {
	t.Parallel()

	st := collectSingleInstanceGuardStatus("cardbot", 1234, func(processName string, selfPID int) (bool, error) {
		return false, errors.New("boom")
	})

	if !st.Enabled {
		t.Fatal("Enabled = false, want true")
	}
	if st.HasOtherProcess {
		t.Fatal("HasOtherProcess = true, want false")
	}
	if st.CheckError == "" {
		t.Fatal("CheckError empty, want value")
	}
}
