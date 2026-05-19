package parse

import (
	"testing"
	"time"
)

func TestParseInputTime(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want time.Duration
	}{
		{
			name: "plain time",
			in:   "01:02:03",
			want: time.Hour + 2*time.Minute + 3*time.Second,
		},
		{
			name: "bracketed time",
			in:   "[14:05:00]",
			want: 14*time.Hour + 5*time.Minute,
		},
		{
			name: "midnight",
			in:   "[00:00:00]",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInputTime(tt.in)
			if err != nil {
				t.Fatalf("ParseInputTime(%q) returned error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ParseInputTime(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseInputTimeInvalid(t *testing.T) {
	if _, err := ParseInputTime("[not-a-time]"); err == nil {
		t.Fatal("ParseInputTime returned nil error for invalid time")
	}
}

func TestFormatOutputTime(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want string
	}{
		{
			name: "zero",
			in:   0,
			want: "[00:00:00]",
		},
		{
			name: "pads components",
			in:   time.Hour + 2*time.Minute + 3*time.Second,
			want: "[01:02:03]",
		},
		{
			name: "supports durations over a day",
			in:   27*time.Hour + 4*time.Minute + 5*time.Second,
			want: "[27:04:05]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatOutputTime(tt.in); got != tt.want {
				t.Fatalf("FormatOutputTime(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
