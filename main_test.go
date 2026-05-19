package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

func TestMainRunsWithExplicitConfigAndEvents(t *testing.T) {
	dir := t.TempDir()

	cfgPath, eventsPath, want := writeMainScenario(t, dir)

	oldArgs := os.Args
	os.Args = []string{"dungeon", cfgPath, eventsPath}
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	got := captureOutput(t, main)
	if got != want {
		t.Fatalf("main output:\n%s\nwant:\n%s", got, want)
	}
}

func TestMainRunsWithDefaultConfigAndEvents(t *testing.T) {
	dir := t.TempDir()
	_, _, want := writeMainScenario(t, dir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	oldArgs := os.Args
	os.Args = []string{"dungeon"}
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	got := captureOutput(t, main)
	if got != want {
		t.Fatalf("main output:\n%s\nwant:\n%s", got, want)
	}
}

func TestMainRunsWithExplicitConfigAndDefaultEvents(t *testing.T) {
	dir := t.TempDir()
	cfgPath, _, want := writeMainScenario(t, dir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	oldArgs := os.Args
	os.Args = []string{"dungeon", cfgPath}
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	got := captureOutput(t, main)
	if got != want {
		t.Fatalf("main output:\n%s\nwant:\n%s", got, want)
	}
}

func TestMainExitsOnTooManyArguments(t *testing.T) {
	if os.Getenv("GO_TEST_MAIN_TOO_MANY_ARGS") == "1" {
		os.Args = []string{"dungeon", "one", "two", "three"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainExitsOnTooManyArguments")
	cmd.Env = append(os.Environ(), "GO_TEST_MAIN_TOO_MANY_ARGS=1")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("main with too many arguments exited successfully, want failure")
	}
	if !strings.Contains(string(out), "Too many arguments.") {
		t.Fatalf("output = %q, want too many arguments message", string(out))
	}
}

func writeMainScenario(t *testing.T, dir string) (cfgPath, eventsPath, want string) {
	t.Helper()

	cfgPath = filepath.Join(dir, "config.json")
	cfg := `{
		"Floors": 2,
		"Monsters": 1,
		"OpenAt": "10:00:00",
		"Duration": 1
	}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	eventsPath = filepath.Join(dir, "events")
	events := "" +
		"[10:00:00] 1 1\n" +
		"[10:01:00] 1 2\n" +
		"[10:02:00] 1 3\n" +
		"[10:03:00] 1 4\n" +
		"[10:04:00] 1 6\n" +
		"[10:05:00] 1 7\n" +
		"[10:06:00] 1 8\n"
	if err := os.WriteFile(eventsPath, []byte(events), 0o600); err != nil {
		t.Fatalf("write events: %v", err)
	}

	want = "" +
		"[10:00:00] Player [1] registered\n" +
		"[10:01:00] Player [1] entered the dungeon\n" +
		"[10:02:00] Player [1] killed the monster\n" +
		"[10:03:00] Player [1] went to the next floor\n" +
		"[10:04:00] Player [1] entered the boss's floor\n" +
		"[10:05:00] Player [1] killed the boss\n" +
		"[10:06:00] Player [1] left the dungeon\n" +
		"\nFinal report:\n" +
		"[SUCCESS] 1 [00:05:00, 00:01:00, 00:02:00] HP:100\n"

	return cfgPath, eventsPath, want
}
