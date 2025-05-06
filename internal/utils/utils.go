package utils

import (
	"biathlonCompetitions/internal/events"
	"fmt"
	"strconv"
	"strings"
)

// SumBoolArray подсчитывает количество истинных (true) значений в массиве.
func SumBoolArray(arr []bool) int {
	total := 0
	for _, b := range arr {
		if b {
			total += 1
		}
	}
	return total
}

// ParseStrToEventData преобразует строку в структурированное событие.
// Формат строки: "[ВРЕМЯ] ID_СОБЫТИЯ ID_УЧАСТНИКА [ДОП.ПАРАМЕТРЫ]"
// Пример: "[12:30:45.000] 5 1234 3"
func ParseStrToEventData(line string) (events.EventData, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return events.EventData{}, fmt.Errorf("not enough elements")
	}

	// Извлечение и нормализация времени вида [HH:MM:SS.ммм]
	eventTimeStr := strings.Trim(parts[0], "[]")

	eventTime, err := ParseStrTimeToTime(eventTimeStr)
	if err != nil {
		return events.EventData{}, fmt.Errorf("incorrect time format: %w", err)
	}

	// Парсинг идентификатора события
	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return events.EventData{}, fmt.Errorf("error parsing the event ID: %v", err)
	}

	// Парсинг идентификатора участника
	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return events.EventData{}, fmt.Errorf(" parsing error the competitor ID: %v", err)
	}

	// Сбор дополнительных параметров (если есть)
	extraParams := ""
	if len(parts) > 3 {
		extraParams = strings.Join(parts[3:], " ")
	}

	return events.EventData{
		ID:           eventID,
		CompetitorID: competitorID,
		Params:       extraParams,
		Time:         eventTime,
	}, nil
}
