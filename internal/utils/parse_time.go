package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Форматы времени для парсинга
const (
	timeFormat1 = "15:04:05.000" // Формат с миллисекундами (HH:MM:SS.ммм)
	timeFormat2 = "15:04:05"     // Формат без миллисекунд (HH:MM:SS)
)

// ParseStrTimeToTime преобразует time.Time.String() в объект time.Time
func ParseStrTimeToTime(strTime string) (time.Time, error) {
	t, err := time.Parse(timeFormat1, strTime)
	if err == nil {
		return t, nil
	}

	return time.Parse(timeFormat2, strTime)
}

// ParseStrTimeToDuration преобразует time.Time.String() в time.Duration
func ParseStrTimeToDuration(strTime string) (time.Duration, error) {
	parts := strings.Split(strTime, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("expected HH:MM:SS format, got: %s", strTime)
	}

	// Парсинг и валидация часов
	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 {
		return 0, fmt.Errorf("invalid hours: %w", err)
	}

	// Парсинг и валидация минут
	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes >= 60 {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	// Парсинг и валидация секунд
	seconds, err := strconv.Atoi(parts[2])
	if err != nil || seconds < 0 || seconds >= 60 {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

// ParseDurationToStrTime преобразует time.Duration в строку формата HH:MM:SS.ммм.
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

// ParseStrDurationToStrTime преобразует time.Duration.String() в time.Time.String()
// Обрабатывает специальные случаи:
//   - "NotStarted" и "NotFinished" возвращаются как есть
//   - Некорректные значения возвращают "InvalidTime"
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
