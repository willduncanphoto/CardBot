package instance

import (
	"errors"
	"os/exec"
	"testing"
)

func TestHasOtherProcessWithRunner_NoMatches(t *testing.T) {
	run := func(name string, args ...string) ([]byte, error) {
		return nil, &exec.ExitError{}
	}

	got, err := hasOtherProcessWithRunner("cardbot", 1234, run)
	if err != nil {
		t.Fatalf("hasOtherProcessWithRunner error: %v", err)
	}
	if got {
		t.Fatal("got true, want false")
	}
}

func TestHasOtherProcessWithRunner_OnlySelf(t *testing.T) {
	run := func(name string, args ...string) ([]byte, error) {
		return []byte("1234\n"), nil
	}

	got, err := hasOtherProcessWithRunner("cardbot", 1234, run)
	if err != nil {
		t.Fatalf("hasOtherProcessWithRunner error: %v", err)
	}
	if got {
		t.Fatal("got true, want false")
	}
}

func TestHasOtherProcessWithRunner_HasAnotherPID(t *testing.T) {
	run := func(name string, args ...string) ([]byte, error) {
		return []byte("1234\n9999\n"), nil
	}

	got, err := hasOtherProcessWithRunner("cardbot", 1234, run)
	if err != nil {
		t.Fatalf("hasOtherProcessWithRunner error: %v", err)
	}
	if !got {
		t.Fatal("got false, want true")
	}
}

func TestHasOtherProcessWithRunner_IgnoresInvalidLines(t *testing.T) {
	run := func(name string, args ...string) ([]byte, error) {
		return []byte("1234\nnot-a-pid\n"), nil
	}

	got, err := hasOtherProcessWithRunner("cardbot", 1234, run)
	if err != nil {
		t.Fatalf("hasOtherProcessWithRunner error: %v", err)
	}
	if got {
		t.Fatal("got true, want false")
	}
}

func TestHasOtherProcessWithRunner_CommandError(t *testing.T) {
	run := func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("boom")
	}

	_, err := hasOtherProcessWithRunner("cardbot", 1234, run)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHasOtherProcessWithRunner_RequiresProcessName(t *testing.T) {
	run := func(name string, args ...string) ([]byte, error) {
		return []byte(""), nil
	}

	_, err := hasOtherProcessWithRunner("", 1234, run)
	if err == nil {
		t.Fatal("expected error")
	}
}
