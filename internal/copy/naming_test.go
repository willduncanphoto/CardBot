package copy

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/illwill/cardbot/internal/analyze"
)

func TestSequenceDigits(t *testing.T) {
	t.Parallel()
	tests := []struct {
		count int
		want  int
	}{
		{0, 3},
		{1, 3},
		{999, 3},
		{1000, 4},
		{9999, 4},
		{10000, 5},
		{250000, 5},
	}

	for _, tt := range tests {
		if got := sequenceDigits(tt.count); got != tt.want {
			t.Errorf("sequenceDigits(%d) = %d, want %d", tt.count, got, tt.want)
		}
	}
}

func TestRenamedRelativePath(t *testing.T) {
	t.Parallel()
	capture := time.Date(2026, 3, 14, 14, 30, 52, 0, time.UTC)
	got := renamedRelativePath("100NIKON/DSC_0001.nef", capture, 12, 4)
	want := filepath.Join("100NIKON", "260314T143052_0012.NEF")
	if got != want {
		t.Fatalf("renamedRelativePath = %q, want %q", got, want)
	}
}

func TestCopy_TimestampNaming(t *testing.T) {
	t.Parallel()
	card := createTestCard(t, map[string]testFileSpec{
		"100NIKON/DSC_0001.NEF": {data: []byte("a"), mtime: date(2026, 3, 8)},
		"100NIKON/DSC_0002.MOV": {data: []byte("b"), mtime: date(2026, 3, 8)},
	})
	dest := t.TempDir()

	ts := time.Date(2026, 3, 14, 14, 30, 52, 0, time.UTC)
	res, err := Run(context.Background(), Options{
		CardPath:   card,
		DestBase:   dest,
		NamingMode: namingModeTimestamp,
		AnalyzeResult: &analyze.Result{
			FileDates: map[string]string{
				"100NIKON/DSC_0001.NEF": "2026-03-14",
				"100NIKON/DSC_0002.MOV": "2026-03-14",
			},
			FileDateTimes: map[string]time.Time{
				"100NIKON/DSC_0001.NEF": ts,
				"100NIKON/DSC_0002.MOV": ts,
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.FilesCopied != 2 {
		t.Fatalf("FilesCopied = %d, want 2", res.FilesCopied)
	}

	assertFileSize(t, filepath.Join(dest, "2026-03-14", "100NIKON", "260314T143052_001.NEF"), 1)
	assertFileSize(t, filepath.Join(dest, "2026-03-14", "100NIKON", "260314T143052_002.MOV"), 1)

	// Original camera names should not be present in timestamp mode.
	if _, err := os.Stat(filepath.Join(dest, "2026-03-14", "100NIKON", "DSC_0001.NEF")); !os.IsNotExist(err) {
		t.Fatal("original filename should not exist in timestamp mode")
	}
}

func TestCopy_TimestampNaming_UsesScanCountForDigits(t *testing.T) {
	t.Parallel()
	card := createTestCard(t, map[string]testFileSpec{
		"100NIKON/DSC_0001.NEF": {data: []byte("a"), mtime: date(2026, 3, 8)},
	})
	dest := t.TempDir()
	ts := time.Date(2026, 3, 14, 14, 30, 52, 0, time.UTC)

	_, err := Run(context.Background(), Options{
		CardPath:   card,
		DestBase:   dest,
		NamingMode: namingModeTimestamp,
		AnalyzeResult: &analyze.Result{
			FileCount: 3048, // simulate card scan count
			FileDates: map[string]string{"100NIKON/DSC_0001.NEF": "2026-03-14"},
			FileDateTimes: map[string]time.Time{
				"100NIKON/DSC_0001.NEF": ts,
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertFileSize(t, filepath.Join(dest, "2026-03-14", "100NIKON", "260314T143052_0001.NEF"), 1)
}
