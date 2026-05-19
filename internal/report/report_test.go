package report

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/l-golofastov/impulse-dungeon-game/internal/models"
	"github.com/l-golofastov/impulse-dungeon-game/internal/runtime"
)

func captureOutput(t *testing.T, f func()) string {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	os.Stdout = writer

	output := make(chan string)
	go func() {
		var builder strings.Builder
		_, _ = io.Copy(&builder, reader)
		output <- builder.String()
	}()

	f()

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	os.Stdout = oldStdout

	out := <-output
	if err := reader.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}

	return out
}

func TestBuildPrintsSortedFinalReport(t *testing.T) {
	game := &runtime.Game{
		Players: map[int]models.Player{
			3: {
				ID:     3,
				Health: 100,
				Dungeon: models.Dungeon{
					Floors: []models.Floor{
						{},
						{BossFloor: true},
					},
				},
			},
			1: {
				ID:     1,
				Health: 35,
				Dungeon: models.Dungeon{
					Completed: true,
					Start:     14*time.Hour + 40*time.Minute,
					Finish:    15*time.Hour + 4*time.Minute,
					Floors: []models.Floor{
						{Cleared: true, TimeToClear: 5 * time.Minute},
						{BossFloor: true, Cleared: true, TimeToClear: 11 * time.Minute},
					},
				},
			},
			2: {
				ID:     2,
				Health: 0,
				Dungeon: models.Dungeon{
					Start:  14*time.Hour + 10*time.Minute,
					Finish: 14*time.Hour + 29*time.Minute,
					Floors: []models.Floor{
						{TimeToClear: 19 * time.Minute},
						{BossFloor: true},
					},
				},
			},
		},
	}

	got := captureOutput(t, NewReport(game).Build)

	want := "" +
		"\nFinal report:\n" +
		"[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35\n" +
		"[FAIL] 2 [00:19:00, 00:00:00, 00:00:00] HP:0\n" +
		"[FAIL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100\n"
	if got != want {
		t.Fatalf("Build output:\n%s\nwant:\n%s", got, want)
	}
}

func TestGetStatus(t *testing.T) {
	report := &Report{}

	tests := []struct {
		name   string
		player models.Player
		want   string
	}{
		{name: "disqualified wins over completed", player: models.Player{Disqualified: true, Dungeon: models.Dungeon{Completed: true}}, want: "[DISQUAL]"},
		{name: "success", player: models.Player{Dungeon: models.Dungeon{Completed: true}}, want: "[SUCCESS]"},
		{name: "fail", player: models.Player{}, want: "[FAIL]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := report.getStatus(tt.player); got != tt.want {
				t.Fatalf("getStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetTimingsAveragesClearedFloorsAndBossTime(t *testing.T) {
	report := &Report{}
	player := models.Player{
		Dungeon: models.Dungeon{
			Start:  10 * time.Hour,
			Finish: 10*time.Hour + 20*time.Minute,
			Floors: []models.Floor{
				{Cleared: true, TimeToClear: 4 * time.Minute},
				{Cleared: true, TimeToClear: 8 * time.Minute},
				{BossFloor: true, Cleared: true, TimeToClear: 6 * time.Minute},
			},
		},
	}

	if got := report.getTimings(player); got != "[00:20:00, 00:06:00, 00:06:00]" {
		t.Fatalf("getTimings() = %q, want %q", got, "[00:20:00, 00:06:00, 00:06:00]")
	}
}

func TestGetTimingsWithNoClearedFloors(t *testing.T) {
	report := &Report{}
	player := models.Player{
		Dungeon: models.Dungeon{
			Start:  10 * time.Hour,
			Finish: 10*time.Hour + 20*time.Minute,
			Floors: []models.Floor{
				{TimeToClear: 4 * time.Minute},
				{BossFloor: true, TimeToClear: 6 * time.Minute},
			},
		},
	}

	if got := report.getTimings(player); got != "[00:20:00, 00:00:00, 00:00:00]" {
		t.Fatalf("getTimings() = %q, want %q", got, "[00:20:00, 00:00:00, 00:00:00]")
	}
}
