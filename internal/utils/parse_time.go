package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat1 = "15:04:05.000"
	timeFormat2 = "15:04:05"
)

func ParseStrTimeToTime(strTime string) (time.Time, error) {
	t, err := time.Parse(timeFormat1, strTime)
	if err == nil {
		return t, nil
	}

	return time.Parse(timeFormat2, strTime)
}

func ParseStrTimeToDuration(strTime string) (time.Duration, error) {
	parts := strings.Split(strTime, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("expected HH:MM:SS format, got: %s", strTime)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 {
		return 0, fmt.Errorf("invalid hours: %w", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes >= 60 {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil || seconds < 0 || seconds >= 60 {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

func ParseDurationToStrTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d",
		hours,
		minutes,
		seconds,
		milliseconds,
	)
}

func ParseStrDurationToStrTime(timeStr string) string {
	if timeStr == "NotStarted" || timeStr == "NotFinished" {
		return timeStr
	}

	dur, err := time.ParseDuration(timeStr)
	if err != nil {
		return "InvalidTime"
	}

	return ParseDurationToStrTime(dur)
}
