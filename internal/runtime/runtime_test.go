package runtime

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/l-golofastov/impulse-dungeon-game/internal/config"
	"github.com/l-golofastov/impulse-dungeon-game/internal/models"
)

func at(hour, minute, second int) time.Duration {
	return time.Duration(hour)*time.Hour +
		time.Duration(minute)*time.Minute +
		time.Duration(second)*time.Second
}

func event(hour, minute, second, playerID, eventID int, extraParam string) models.Event {
	return models.Event{
		Time:       at(hour, minute, second),
		PlayerID:   playerID,
		EventID:    eventID,
		ExtraParam: extraParam,
	}
}

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

func TestRunProcessesReadmeScenario(t *testing.T) {
	cfg := &config.RuntimeConfig{
		Floors:   2,
		Monsters: 2,
		OpenAt:   at(14, 5, 0),
		CloseAt:  at(16, 5, 0),
	}

	events := []models.Event{
		event(14, 0, 0, 1, 1, ""),
		event(14, 0, 0, 2, 1, ""),
		event(14, 10, 0, 2, 2, ""),
		event(14, 10, 0, 3, 2, ""),
		event(14, 11, 0, 2, 5, ""),
		event(14, 12, 0, 3, 3, ""),
		event(14, 14, 0, 2, 3, ""),
		event(14, 27, 0, 2, 11, "60"),
		event(14, 29, 0, 2, 11, "50"),
		event(14, 40, 0, 1, 2, ""),
		event(14, 41, 0, 1, 3, ""),
		event(14, 44, 0, 1, 11, "50"),
		event(14, 45, 0, 1, 3, ""),
		event(14, 48, 0, 1, 4, ""),
		event(14, 48, 0, 1, 6, ""),
		event(14, 49, 0, 1, 11, "25"),
		event(14, 49, 2, 1, 10, "80"),
		event(14, 50, 0, 1, 11, "65"),
		event(14, 59, 0, 1, 7, ""),
		event(15, 4, 0, 1, 8, ""),
	}

	game := NewGame(cfg, events)
	got := captureOutput(t, game.Run)

	want := "" +
		"[14:00:00] Player [1] registered\n" +
		"[14:00:00] Player [2] registered\n" +
		"[14:10:00] Player [2] entered the dungeon\n" +
		"[14:10:00] Player [3] is disqualified\n" +
		"[14:11:00] Player [2] makes imposible move [5]\n" +
		"[14:14:00] Player [2] killed the monster\n" +
		"[14:27:00] Player [2] recieved [60] of damage\n" +
		"[14:29:00] Player [2] recieved [50] of damage\n" +
		"[14:29:00] Player [2] is dead\n" +
		"[14:40:00] Player [1] entered the dungeon\n" +
		"[14:41:00] Player [1] killed the monster\n" +
		"[14:44:00] Player [1] recieved [50] of damage\n" +
		"[14:45:00] Player [1] killed the monster\n" +
		"[14:48:00] Player [1] went to the next floor\n" +
		"[14:48:00] Player [1] entered the boss's floor\n" +
		"[14:49:00] Player [1] recieved [25] of damage\n" +
		"[14:49:02] Player [1] has restored [80] of health\n" +
		"[14:50:00] Player [1] recieved [65] of damage\n" +
		"[14:59:00] Player [1] killed the boss\n" +
		"[15:04:00] Player [1] left the dungeon\n"
	if got != want {
		t.Fatalf("Run output:\n%s\nwant:\n%s", got, want)
	}

	player1 := game.Players[1]
	if !player1.Completed || !player1.LeftDungeon {
		t.Fatalf("player 1 completion/left flags = %v/%v, want true/true", player1.Completed, player1.LeftDungeon)
	}
	if player1.Health != 35 {
		t.Fatalf("player 1 health = %d, want 35", player1.Health)
	}
	if player1.Finish-player1.Start != 24*time.Minute {
		t.Fatalf("player 1 full time = %v, want 24m", player1.Finish-player1.Start)
	}
	if player1.Floors[0].TimeToClear != 5*time.Minute {
		t.Fatalf("player 1 first floor time = %v, want 5m", player1.Floors[0].TimeToClear)
	}
	if player1.Floors[1].TimeToClear != 11*time.Minute {
		t.Fatalf("player 1 boss time = %v, want 11m", player1.Floors[1].TimeToClear)
	}

	player2 := game.Players[2]
	if player2.Health != 0 {
		t.Fatalf("player 2 health = %d, want 0", player2.Health)
	}
	if player2.Finish != at(14, 29, 0) {
		t.Fatalf("player 2 finish = %v, want 14:29:00", player2.Finish)
	}

	player3 := game.Players[3]
	if !player3.Disqualified {
		t.Fatal("player 3 should be disqualified")
	}
}

func TestRunHandlesImpossibleMovesBeforeDungeonAndUnknownEvent(t *testing.T) {
	cfg := &config.RuntimeConfig{
		Floors:   3,
		Monsters: 1,
		OpenAt:   at(10, 0, 0),
		CloseAt:  at(11, 0, 0),
	}
	events := []models.Event{
		event(9, 59, 0, 1, 2, ""),
		event(10, 0, 0, 1, 3, ""),
		event(10, 1, 0, 1, 4, ""),
		event(10, 2, 0, 1, 5, ""),
		event(10, 3, 0, 1, 6, ""),
		event(10, 4, 0, 1, 8, ""),
		event(10, 5, 0, 1, 10, "10"),
		event(10, 6, 0, 1, 11, "10"),
		event(10, 7, 0, 1, 99, ""),
	}

	game := NewGame(cfg, events)
	got := captureOutput(t, game.Run)

	want := "" +
		"[09:59:00] Player [1] makes imposible move [2]\n" +
		"[10:00:00] Player [1] makes imposible move [3]\n" +
		"[10:01:00] Player [1] makes imposible move [4]\n" +
		"[10:02:00] Player [1] makes imposible move [5]\n" +
		"[10:03:00] Player [1] makes imposible move [6]\n" +
		"[10:04:00] Player [1] makes imposible move [8]\n" +
		"[10:05:00] Player [1] makes imposible move [10]\n" +
		"[10:06:00] Player [1] makes imposible move [11]\n" +
		"Unknown Event\n"
	if got != want {
		t.Fatalf("Run output:\n%s\nwant:\n%s", got, want)
	}
}

func TestRunHandlesMovementAndRepeatedClearedActions(t *testing.T) {
	cfg := &config.RuntimeConfig{
		Floors:   3,
		Monsters: 1,
		OpenAt:   at(10, 0, 0),
		CloseAt:  at(12, 0, 0),
	}
	events := []models.Event{
		event(10, 0, 0, 1, 1, ""),
		event(10, 0, 30, 1, 1, ""),
		event(10, 1, 0, 1, 2, ""),
		event(10, 2, 0, 1, 10, "bad"),
		event(10, 3, 0, 1, 11, "bad"),
		event(10, 4, 0, 1, 4, ""),
		event(10, 5, 0, 1, 3, ""),
		event(10, 6, 0, 1, 3, ""),
		event(10, 7, 0, 1, 4, ""),
		event(10, 8, 0, 1, 7, ""),
		event(10, 9, 0, 1, 5, ""),
		event(10, 10, 0, 1, 5, ""),
		event(10, 11, 0, 1, 4, ""),
		event(10, 12, 0, 1, 3, ""),
		event(10, 13, 0, 1, 4, ""),
		event(10, 14, 0, 1, 6, ""),
		event(10, 15, 0, 1, 7, ""),
		event(10, 16, 0, 1, 7, ""),
		event(10, 17, 0, 1, 4, ""),
		event(10, 18, 0, 1, 8, ""),
		event(10, 19, 0, 1, 3, ""),
	}

	game := NewGame(cfg, events)
	got := captureOutput(t, game.Run)

	want := "" +
		"[10:00:00] Player [1] registered\n" +
		"[10:01:00] Player [1] entered the dungeon\n" +
		"[10:02:00] Player [1] makes imposible move [10]\n" +
		"[10:03:00] Player [1] makes imposible move [11]\n" +
		"[10:04:00] Player [1] makes imposible move [4]\n" +
		"[10:05:00] Player [1] killed the monster\n" +
		"[10:06:00] Player [1] makes imposible move [3]\n" +
		"[10:07:00] Player [1] went to the next floor\n" +
		"[10:08:00] Player [1] makes imposible move [7]\n" +
		"[10:09:00] Player [1] went to the previous floor\n" +
		"[10:10:00] Player [1] makes imposible move [5]\n" +
		"[10:11:00] Player [1] went to the next floor\n" +
		"[10:12:00] Player [1] killed the monster\n" +
		"[10:13:00] Player [1] went to the next floor\n" +
		"[10:14:00] Player [1] entered the boss's floor\n" +
		"[10:15:00] Player [1] killed the boss\n" +
		"[10:16:00] Player [1] makes imposible move [7]\n" +
		"[10:17:00] Player [1] makes imposible move [4]\n" +
		"[10:18:00] Player [1] left the dungeon\n"
	if got != want {
		t.Fatalf("Run output:\n%s\nwant:\n%s", got, want)
	}

	player := game.Players[1]
	if !player.Completed || !player.LeftDungeon {
		t.Fatalf("completed/left = %v/%v, want true/true", player.Completed, player.LeftDungeon)
	}
	if player.Floors[0].TimeToClear != 4*time.Minute {
		t.Fatalf("first floor time = %v, want 4m", player.Floors[0].TimeToClear)
	}
	if player.Floors[1].TimeToClear != 3*time.Minute {
		t.Fatalf("second floor time = %v, want 3m", player.Floors[1].TimeToClear)
	}
	if player.Floors[2].TimeToClear != 2*time.Minute {
		t.Fatalf("boss floor time = %v, want 2m", player.Floors[2].TimeToClear)
	}
}

func TestRunHandlesHealthCannotContinueAndClosedDungeon(t *testing.T) {
	cfg := &config.RuntimeConfig{
		Floors:   2,
		Monsters: 1,
		OpenAt:   at(10, 0, 0),
		CloseAt:  at(10, 30, 0),
	}
	events := []models.Event{
		event(10, 0, 0, 1, 1, ""),
		event(10, 1, 0, 1, 2, ""),
		event(10, 2, 0, 1, 11, "90"),
		event(10, 3, 0, 1, 10, "200"),
		event(10, 4, 0, 1, 9, "broken sword"),
		event(10, 5, 0, 1, 11, "10"),
		event(10, 6, 0, 2, 9, "no torch"),
		event(10, 7, 0, 2, 1, ""),
		event(10, 8, 0, 3, 1, ""),
		event(10, 31, 0, 3, 2, ""),
	}

	game := NewGame(cfg, events)
	got := captureOutput(t, game.Run)

	want := "" +
		"[10:00:00] Player [1] registered\n" +
		"[10:01:00] Player [1] entered the dungeon\n" +
		"[10:02:00] Player [1] recieved [90] of damage\n" +
		"[10:03:00] Player [1] has restored [200] of health\n" +
		"[10:04:00] Player [1] cannot continue due to [broken sword]\n" +
		"[10:06:00] Player [2] cannot continue due to [no torch]\n" +
		"[10:08:00] Player [3] registered\n"
	if got != want {
		t.Fatalf("Run output:\n%s\nwant:\n%s", got, want)
	}

	player1 := game.Players[1]
	if !player1.Disqualified {
		t.Fatal("player 1 should be disqualified")
	}
	if player1.Health != 100 {
		t.Fatalf("player 1 health = %d, want capped health 100", player1.Health)
	}
	if player1.Finish != at(10, 4, 0) {
		t.Fatalf("player 1 finish = %v, want 10:04:00", player1.Finish)
	}

	player2 := game.Players[2]
	if !player2.Disqualified {
		t.Fatal("player 2 should be disqualified")
	}
	if player2.Finish != 0 {
		t.Fatalf("player 2 finish = %v, want zero because dungeon was not started", player2.Finish)
	}

	player3 := game.Players[3]
	if player3.CurrentFloor != 0 {
		t.Fatalf("player 3 current floor = %d, want 0 because event after close is ignored", player3.CurrentFloor)
	}
}

func TestNewGameStoresConfigEventsAndInitializesPlayers(t *testing.T) {
	cfg := &config.RuntimeConfig{Floors: 2, Monsters: 1}
	events := []models.Event{event(1, 0, 0, 10, 1, "")}

	game := NewGame(cfg, events)

	if game.Config != cfg {
		t.Fatal("NewGame did not keep config pointer")
	}
	if len(game.Events) != 1 || game.Events[0] != events[0] {
		t.Fatalf("NewGame events = %+v, want %+v", game.Events, events)
	}
	if game.Players == nil {
		t.Fatal("NewGame Players map is nil")
	}
}
