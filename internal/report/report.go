package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/l-golofastov/impulse-dungeon-game/internal/models"
	"github.com/l-golofastov/impulse-dungeon-game/internal/parse"
	"github.com/l-golofastov/impulse-dungeon-game/internal/runtime"
)

type Report struct {
	game *runtime.Game
}

func NewReport(game *runtime.Game) *Report {
	return &Report{game: game}
}

func (r *Report) Build() {
	keys := make([]int, 0, len(r.game.Players))

	for k := range r.game.Players {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Printf("\nFinal report:\n")

	for _, k := range keys {
		player := r.game.Players[k]

		status := r.getStatus(player)
		timings := r.getTimings(player)

		fmt.Printf("%s %d %s HP:%d\n", status, player.ID, timings, player.Health)
	}
}

func (r *Report) getStatus(player models.Player) string {
	var status string

	if player.Disqualified {
		status = "[DISQUAL]"
	} else if player.Completed {
		status = "[SUCCESS]"
	} else {
		status = "[FAIL]"
	}

	return status
}

func (r *Report) getTimings(player models.Player) string {
	fullTime := parse.FormatOutputTime(player.Finish - player.Start)
	fullTime = strings.Trim(fullTime, "[]")

	avgTime := r.getAvgTime(player)
	avgTime = strings.Trim(avgTime, "[]")

	bossTime := r.getBossTime(player)
	bossTime = strings.Trim(bossTime, "[]")

	return fmt.Sprintf("[%s, %s, %s]", fullTime, avgTime, bossTime)
}

func (r *Report) getBossTime(player models.Player) string {
	var time string

	bossFloor := player.Floors[len(player.Floors)-1]
	if bossFloor.Cleared {
		time = parse.FormatOutputTime(bossFloor.TimeToClear)
	} else {
		time = parse.FormatOutputTime(0)
	}

	return time
}

func (r *Report) getAvgTime(player models.Player) string {
	var duration time.Duration = 0
	clearedFloors := 0

	for i := 0; i < len(player.Floors)-1; i++ {
		floor := player.Floors[i]
		if floor.Cleared {
			duration += floor.TimeToClear
			clearedFloors++
		}
	}

	if clearedFloors > 0 {
		duration /= time.Duration(clearedFloors)
	}

	return parse.FormatOutputTime(duration)
}
