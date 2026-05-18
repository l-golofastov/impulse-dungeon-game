package models

import "time"

type Player struct {
	ID           int
	Health       int
	Registered   bool
	Disqualified bool
	LeftDungeon  bool
	Dungeon
}

type Event struct {
	Time       time.Duration
	PlayerID   int
	EventID    int
	ExtraParam string
}

type Floor struct {
	LastEventTime time.Duration
	MonstersAlive int
	BossFloor     bool
	TimeToClear   time.Duration
	Cleared       bool
}

type Dungeon struct {
	Start        time.Duration
	Finish       time.Duration
	CurrentFloor int
	Floors       []Floor
	Completed    bool
}
