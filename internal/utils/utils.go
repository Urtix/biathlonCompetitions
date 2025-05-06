package utils

import (
	"biathlonCompetitions/internal/events"
	"fmt"
	"strconv"
	"strings"
)

func SumIntArray(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func ParseStrToEventData(line string) (events.EventData, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return events.EventData{}, fmt.Errorf("недостаточно элементов")
	}

	eventTimeStr := strings.Trim(parts[0], "[]")

	eventTime, err := ParseStrTimeToTime(eventTimeStr)
	if err != nil {
		return events.EventData{}, fmt.Errorf("некорректный формат времени: %w", err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return events.EventData{}, fmt.Errorf("ошибка парсинга iD события: %v", err)
	}

	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return events.EventData{}, fmt.Errorf("ошибка парсинга iD участника: %v", err)
	}

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
