package parse

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseEvents(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events")
	content := "" +
		"[14:00:00] 1 1\n" +
		"[14:10:05] 2 9 out of arrows\n" +
		"[14:11:00] 2 10 25\n"

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write events file: %v", err)
	}

	events, err := ParseEvents(path)
	if err != nil {
		t.Fatalf("ParseEvents returned error: %v", err)
	}

	if len(events) != 3 {
		t.Fatalf("ParseEvents returned %d events, want 3", len(events))
	}

	if events[0].Time != 14*time.Hour || events[0].PlayerID != 1 || events[0].EventID != 1 || events[0].ExtraParam != "" {
		t.Fatalf("unexpected first event: %+v", events[0])
	}

	wantSecondTime := 14*time.Hour + 10*time.Minute + 5*time.Second
	if events[1].Time != wantSecondTime || events[1].PlayerID != 2 || events[1].EventID != 9 || events[1].ExtraParam != "out of arrows" {
		t.Fatalf("unexpected second event with multi-word extra param: %+v", events[1])
	}

	if events[2].ExtraParam != "25" {
		t.Fatalf("third event ExtraParam = %q, want %q", events[2].ExtraParam, "25")
	}
}

func TestParseEventsMissingFile(t *testing.T) {
	if _, err := ParseEvents(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("ParseEvents returned nil error for missing file")
	}
}

func TestParseEventsInvalidLine(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{name: "too few fields", line: "[14:00:00] 1"},
		{name: "invalid time", line: "[99:00:00] 1 1"},
		{name: "invalid player id", line: "[14:00:00] player 1"},
		{name: "invalid event id", line: "[14:00:00] 1 event"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "events")
			if err := os.WriteFile(path, []byte(tt.line+"\n"), 0o600); err != nil {
				t.Fatalf("write events file: %v", err)
			}

			_, err := ParseEvents(path)
			if err == nil {
				t.Fatal("ParseEvents returned nil error for invalid line")
			}
		})
	}
}

func TestParseSingleInputEventErrors(t *testing.T) {
	tests := []struct {
		name string
		line string
		want error
	}{
		{name: "too few fields", line: "[14:00:00] 1", want: ErrInvalidInputEvent},
		{name: "invalid player id", line: "[14:00:00] player 1", want: ErrInvalidInputEvent},
		{name: "invalid event id", line: "[14:00:00] 1 event", want: ErrInvalidInputEvent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSingleInputEvent(tt.line)
			if !errors.Is(err, tt.want) {
				t.Fatalf("parseSingleInputEvent(%q) error = %v, want %v", tt.line, err, tt.want)
			}
		})
	}
}
