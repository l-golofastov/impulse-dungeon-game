package parse

import (
	"fmt"
	"strings"
	"time"
)

func ParseInputTime(s string) (time.Duration, error) {
	s = strings.Trim(s, "[]")

	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}

	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second, nil
}

func FormatOutputTime(t time.Duration) string {
	total := int(t.Seconds())

	h := total / 3600
	m := total % 3600 / 60
	s := total % 60

	return fmt.Sprintf("[%02d:%02d:%02d]", h, m, s)
}
