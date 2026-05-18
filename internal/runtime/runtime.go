package runtime

import (
	"fmt"
	"strconv"

	"github.com/l-golofastov/impulse-dungeon-game/internal/config"
	"github.com/l-golofastov/impulse-dungeon-game/internal/models"
	"github.com/l-golofastov/impulse-dungeon-game/internal/parse"
)

type Game struct {
	Config  *config.RuntimeConfig
	Events  []models.Event
	Players map[int]models.Player
}

func NewGame(cfg *config.RuntimeConfig, events []models.Event) *Game {
	players := make(map[int]models.Player)
	return &Game{
		Config:  cfg,
		Events:  events,
		Players: players,
	}
}

func (g *Game) Run() {
	for _, event := range g.Events {
		if _, ok := g.Players[event.PlayerID]; !ok {
			g.createPlayer(event.PlayerID)
		}

		g.processEvent(event)
	}
}

func (g *Game) ifChallengeEnded(event models.Event) bool {
	player := g.Players[event.PlayerID]
	if event.Time > g.Config.CloseAt || player.LeftDungeon || player.Disqualified || player.Health == 0 {
		return true
	}

	return false
}

func (g *Game) processEvent(event models.Event) {
	if g.ifChallengeEnded(event) {
		return
	}

	switch event.EventID {
	case 1:
		g.processEvent1(event)
	case 2:
		g.processEvent2(event)
	case 3:
		g.processEvent3(event)
	case 4:
		g.processEvent4(event)
	case 5:
		g.processEvent5(event)
	case 6:
		g.processEvent6(event)
	case 7:
		g.processEvent7(event)
	case 8:
		g.processEvent8(event)
	case 9:
		g.processEvent9(event)
	case 10:
		g.processEvent10(event)
	case 11:
		g.processEvent11(event)
	default:
		fmt.Println("Unknown Event")
	}
}

func (g *Game) processImpossibleMove(event models.Event) {
	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] makes imposible move [%d]\n", time, event.PlayerID, event.EventID)
}

func (g *Game) addTime(event models.Event) {
	player := g.Players[event.PlayerID]
	currFloor := player.Floors[player.CurrentFloor-1]

	lastEnterTime := currFloor.LastEventTime
	valuableTime := event.Time - lastEnterTime

	currFloor.TimeToClear += valuableTime
	currFloor.LastEventTime = event.Time

	player.Floors[player.CurrentFloor-1] = currFloor
	g.Players[event.PlayerID] = player
}

func (g *Game) processEvent1(event models.Event) {
	player := g.Players[event.PlayerID]
	if !player.Registered {
		player.Registered = true
		g.Players[event.PlayerID] = player

		time := parse.FormatOutputTime(event.Time)
		fmt.Printf("%s Player [%d] registered\n", time, event.PlayerID)
	}
}

func (g *Game) processEvent2(event models.Event) {
	player := g.Players[event.PlayerID]
	if event.Time < g.Config.OpenAt {
		g.processImpossibleMove(event)
	} else if player.Registered {
		player.CurrentFloor++
		player.Start = event.Time
		player.Floors[0].LastEventTime = event.Time

		g.Players[event.PlayerID] = player

		time := parse.FormatOutputTime(event.Time)
		fmt.Printf("%s Player [%d] entered the dungeon\n", time, event.PlayerID)
	} else {
		player.Disqualified = true
		g.Players[event.PlayerID] = player

		time := parse.FormatOutputTime(event.Time)
		fmt.Printf("%s Player [%d] is disqualified\n", time, event.PlayerID)
	}
}

func (g *Game) processEvent3(event models.Event) {
	player := g.Players[event.PlayerID]
	bossFloorId := g.Config.Floors

	if player.CurrentFloor == 0 || player.CurrentFloor == bossFloorId { //not entered the dungeon or current floor is boss floor
		g.processImpossibleMove(event)
		return
	}

	currFloor := player.Floors[player.CurrentFloor-1]

	if currFloor.Cleared {
		g.processImpossibleMove(event)
		return
	}

	currFloor.MonstersAlive--

	if currFloor.MonstersAlive == 0 {
		currFloor.Cleared = true
	}

	player.Floors[player.CurrentFloor-1] = currFloor

	g.addTime(event)

	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] killed the monster\n", time, event.PlayerID)
}

func (g *Game) processEvent4(event models.Event) {
	player := g.Players[event.PlayerID]
	bossFloorId := g.Config.Floors

	if player.CurrentFloor == 0 || player.CurrentFloor == bossFloorId { //not entered the dungeon or current floor is boss floor
		g.processImpossibleMove(event)
		return
	}

	currFloor := player.Floors[player.CurrentFloor-1]

	if !currFloor.Cleared {
		g.processImpossibleMove(event)
		return
	}

	player.CurrentFloor++

	currFloor = player.Floors[player.CurrentFloor-1]
	if !currFloor.Cleared {
		currFloor.LastEventTime = event.Time
	}

	player.Floors[player.CurrentFloor-1] = currFloor
	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] went to the next floor\n", time, event.PlayerID)
}

func (g *Game) processEvent5(event models.Event) {
	player := g.Players[event.PlayerID]

	if player.CurrentFloor == 0 || player.CurrentFloor == 1 { //not entered the dungeon or current floor 1st floor
		g.processImpossibleMove(event)
		return
	}

	if !player.Floors[player.CurrentFloor-1].Cleared {
		g.addTime(event)
	}

	player.CurrentFloor--
	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] went to the previous floor\n", time, event.PlayerID)
}

func (g *Game) processEvent6(event models.Event) {
	player := g.Players[event.PlayerID]
	bossFloorId := g.Config.Floors

	if player.CurrentFloor != bossFloorId {
		g.processImpossibleMove(event)
		return
	}

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] entered the boss's floor\n", time, event.PlayerID)
}

func (g *Game) processEvent7(event models.Event) {
	player := g.Players[event.PlayerID]
	currFloor := player.Floors[player.CurrentFloor-1]

	if !currFloor.BossFloor || currFloor.Cleared {
		g.processImpossibleMove(event)
		return
	}

	currFloor.Cleared = true
	player.Floors[player.CurrentFloor-1] = currFloor

	player.Completed = true
	g.addTime(event)

	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] killed the boss\n", time, event.PlayerID)
}

func (g *Game) processEvent8(event models.Event) {
	player := g.Players[event.PlayerID]

	if player.CurrentFloor == 0 {
		g.processImpossibleMove(event)
		return
	}

	player.LeftDungeon = true
	player.Finish = event.Time
	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] left the dungeon\n", time, event.PlayerID)
}

func (g *Game) processEvent9(event models.Event) {
	player := g.Players[event.PlayerID]

	player.Disqualified = true
	if player.Start != 0 {
		player.Finish = event.Time
	}
	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] cannot continue due to [%s]\n", time, event.PlayerID, event.ExtraParam)
}

func (g *Game) processEvent10(event models.Event) {
	player := g.Players[event.PlayerID]

	if player.CurrentFloor == 0 {
		g.processImpossibleMove(event)
		return
	}

	hp, err := strconv.Atoi(event.ExtraParam)
	if err != nil {
		g.processImpossibleMove(event)
		return
	}

	newHp := hp + player.Health
	if newHp > 100 {
		newHp = 100
	}
	player.Health = newHp
	g.addTime(event)

	g.Players[event.PlayerID] = player

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] has restored [%d] of health\n", time, event.PlayerID, hp)
}

func (g *Game) processEvent11(event models.Event) {
	player := g.Players[event.PlayerID]

	if player.CurrentFloor == 0 {
		g.processImpossibleMove(event)
		return
	}

	hp, err := strconv.Atoi(event.ExtraParam)
	if err != nil {
		g.processImpossibleMove(event)
		return
	}

	time := parse.FormatOutputTime(event.Time)
	fmt.Printf("%s Player [%d] recieved [%d] of damage\n", time, event.PlayerID, hp)

	newHp := player.Health - hp
	if newHp <= 0 {
		newHp = 0
		player.Finish = event.Time
		fmt.Printf("%s Player [%d] is dead\n", time, event.PlayerID)
	}

	player.Health = newHp
	g.addTime(event)

	g.Players[event.PlayerID] = player
}

func (g *Game) createPlayer(id int) {
	dungeon := g.createDungeon()

	player := models.Player{
		ID:      id,
		Health:  100,
		Dungeon: dungeon,
	}

	g.Players[id] = player
}

func (g *Game) createDungeon() models.Dungeon {
	floors := g.createFloors()

	dungeon := models.Dungeon{
		Floors: floors,
	}

	return dungeon
}

func (g *Game) createFloors() []models.Floor {
	floors := make([]models.Floor, g.Config.Floors)

	for i := 0; i < g.Config.Floors-1; i++ {
		floors[i] = models.Floor{
			MonstersAlive: g.Config.Monsters,
		}
	}

	bossFloor := g.createBossFloor()
	floors[len(floors)-1] = bossFloor

	return floors
}

func (g *Game) createBossFloor() models.Floor {
	bossFloor := models.Floor{
		MonstersAlive: 0,
		BossFloor:     true,
	}

	return bossFloor
}
