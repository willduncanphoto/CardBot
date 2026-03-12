// Package dotfile reads and writes .cardbot files on memory cards.
package dotfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const fileName = ".cardbot"

// Status represents the copy state of a card.
type Status struct {
	Copied     bool      // Whether the card has been copied
	CopiedAt   time.Time // When the last copy completed
	CopiedDest string    // Where files were copied to
}

// dotfileSchema is the on-disk JSON structure (read-only for now).
type dotfileSchema struct {
	Schema      string `json:"$schema"`
	LastCopied  string `json:"last_copied"`
	Destination string `json:"destination"`
}

// Read checks for a .cardbot file on the card and returns its status.
// Returns a "New" status (Copied=false) if the file doesn't exist or can't be parsed.
func Read(cardPath string) Status {
	data, err := os.ReadFile(filepath.Join(cardPath, fileName))
	if err != nil {
		return Status{}
	}

	var df dotfileSchema
	if err := json.Unmarshal(data, &df); err != nil {
		return Status{}
	}

	t, err := time.Parse(time.RFC3339, df.LastCopied)
	if err != nil {
		return Status{}
	}

	return Status{
		Copied:     true,
		CopiedAt:   t,
		CopiedDest: df.Destination,
	}
}

// FormatStatus returns a display string for the card status.
//
//	"New"
//	"Copied on 2026-03-08"
func FormatStatus(s Status) string {
	if !s.Copied {
		return "New"
	}
	return "Copied on " + s.CopiedAt.Format("2006-01-02 15:04")
}
