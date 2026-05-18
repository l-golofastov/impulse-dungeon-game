package main

import (
	"fmt"
	"log"
	"os"

	"github.com/l-golofastov/impulse-dungeon-game/internal/config"
	"github.com/l-golofastov/impulse-dungeon-game/internal/parse"
	rep "github.com/l-golofastov/impulse-dungeon-game/internal/report"
	"github.com/l-golofastov/impulse-dungeon-game/internal/runtime"
)

func main() {
	cfgPath, eventsPath := "config.json", "events"

	if len(os.Args) == 3 {
		cfgPath = os.Args[1]
		eventsPath = os.Args[2]
	} else if len(os.Args) == 2 {
		cfgPath = os.Args[1]
	} else if len(os.Args) > 3 {
		fmt.Println("Too many arguments.")
		fmt.Println("Usage: ./<executable> [config_path] [events_path]")
		os.Exit(1)
	}

	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	events, err := parse.ParseEvents(eventsPath)
	if err != nil {
		log.Fatal(err)
	}

	game := runtime.NewGame(cfg, events)
	game.Run()

	report := rep.NewReport(game)
	report.Build()
}
