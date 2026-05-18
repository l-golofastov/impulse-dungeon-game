package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/l-golofastov/impulse-dungeon-game/internal/parse"
)

type jsonConfig struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

type RuntimeConfig struct {
	Floors   int
	Monsters int
	OpenAt   time.Duration
	CloseAt  time.Duration
}

func ReadConfig(path string) (*RuntimeConfig, error) {
	jsonCfg, err := readJsonConfig(path)
	if err != nil {
		return nil, err
	}

	openAt, err := parse.ParseInputTime(jsonCfg.OpenAt)
	if err != nil {
		return nil, err
	}

	closeAt := openAt + time.Duration(jsonCfg.Duration)*time.Hour

	cfg := RuntimeConfig{
		Floors:   jsonCfg.Floors,
		Monsters: jsonCfg.Monsters,
		OpenAt:   openAt,
		CloseAt:  closeAt,
	}

	return &cfg, nil
}

func readJsonConfig(path string) (*jsonConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg jsonConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
