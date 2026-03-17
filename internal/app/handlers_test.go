package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/illwill/cardbot/internal/config"
	cardcopy "github.com/illwill/cardbot/internal/copy"
	"github.com/illwill/cardbot/internal/detect"
)

func TestHandleCopyCmd_NotReady(t *testing.T) {
	cardPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(cardPath, "DCIM"), 0o755); err != nil {
		t.Fatalf("mkdir dcim: %v", err)
	}

	cfg := config.Defaults()
	cfg.Destination.Path = t.TempDir()
	fd := newFakeDetector()

	called := 0
	a := New(Config{
		Cfg:         cfg,
		DryRun:      true,
		newDetector: func() cardDetector { return fd },
		runCopy: func(ctx context.Context, opts cardcopy.Options, onProgress cardcopy.ProgressFunc) (*cardcopy.Result, error) {
			called++
			return &cardcopy.Result{}, nil
		},
	})
	a.detector = fd

	card := &detect.Card{Path: cardPath, Name: "CARD"}
	a.currentCard = card
	a.phase = phaseAnalyzing

	out := captureStdout(t, func() {
		a.handleCopyCmd(card, "all")
	})

	if called != 0 {
		t.Fatalf("copy runner called %d times, want 0", called)
	}
	if !strings.Contains(out, "Still scanning card. Please wait.") {
		t.Fatalf("expected readiness warning, got:\n%s", out)
	}
}

func TestHandleCopyCmd_Allowed(t *testing.T) {
	cardPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(cardPath, "DCIM"), 0o755); err != nil {
		t.Fatalf("mkdir dcim: %v", err)
	}

	cfg := config.Defaults()
	cfg.Destination.Path = t.TempDir()
	fd := newFakeDetector()

	called := 0
	a := New(Config{
		Cfg:         cfg,
		DryRun:      true,
		newDetector: func() cardDetector { return fd },
		runCopy: func(ctx context.Context, opts cardcopy.Options, onProgress cardcopy.ProgressFunc) (*cardcopy.Result, error) {
			called++
			return &cardcopy.Result{}, nil
		},
	})
	a.detector = fd

	card := &detect.Card{Path: cardPath, Name: "CARD"}
	a.currentCard = card
	a.phase = phaseReady
	a.copiedModes = make(map[string]bool)

	a.handleCopyCmd(card, "all")

	if called != 1 {
		t.Fatalf("copy runner called %d times, want 1", called)
	}
}

func TestIsTracked(t *testing.T) {
	t.Parallel()

	a := &App{
		currentCard: &detect.Card{Path: "/card/current"},
		cardQueue: []*detect.Card{
			{Path: "/card/queued1"},
			{Path: "/card/queued2"},
		},
	}

	if !a.isTracked("/card/current") {
		t.Fatal("expected current card to be tracked")
	}
	if !a.isTracked("/card/queued2") {
		t.Fatal("expected queued card to be tracked")
	}
	if a.isTracked("/card/unknown") {
		t.Fatal("did not expect unknown card to be tracked")
	}
}

func TestResumeScanningIfIdle(t *testing.T) {
	cfg := config.Defaults()
	a := New(Config{Cfg: cfg})

	a.resumeScanningIfIdle()
	if a.spinner == nil {
		t.Fatal("expected spinner when idle")
	}
	a.stopScanning()

	a.currentCard = &detect.Card{Path: "/card"}
	a.resumeScanningIfIdle()
	if a.spinner != nil {
		t.Fatal("did not expect spinner when card is active")
	}
}
