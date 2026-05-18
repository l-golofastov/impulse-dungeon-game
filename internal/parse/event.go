package parse

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/l-golofastov/impulse-dungeon-game/internal/models"
)

func ParseEvents(path string) ([]models.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	events := make([]models.Event, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		eventStr := scanner.Text()
		event, err := parseSingleInputEvent(eventStr)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func parseSingleInputEvent(eventStr string) (models.Event, error) {
	split := strings.SplitN(eventStr, " ", 4)

	if len(split) < 3 {
		return models.Event{}, ErrInvalidInputEvent
	} else if len(split) == 3 {
		event, err := newInputEventFromSlice(split, false)
		if err != nil {
			return models.Event{}, err
		}

		return event, nil
	} else if len(split) == 4 {
		event, err := newInputEventFromSlice(split, true)
		if err != nil {
			return models.Event{}, err
		}

		return event, nil
	}

	return models.Event{}, ErrInvalidInputEvent
}

func newInputEventFromSlice(split []string, extraParamIncluded bool) (models.Event, error) {
	timeStr := split[0]
	time, err := ParseInputTime(timeStr)
	if err != nil {
		return models.Event{}, err
	}

	playerIDStr := split[1]
	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		return models.Event{}, ErrInvalidInputEvent
	}

	eventIDStr := split[2]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		return models.Event{}, ErrInvalidInputEvent
	}

	var extraParam string
	if extraParamIncluded {
		extraParam = split[3]
	}

	event := models.Event{
		Time:       time,
		PlayerID:   playerID,
		EventID:    eventID,
		ExtraParam: extraParam,
	}

	return event, nil
}
