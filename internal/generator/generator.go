package generator

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/events"
	"biathlonCompetitions/internal/observer"
	"biathlonCompetitions/internal/utils"
	"fmt"
	"time"
)

func NewEventGenerator(config config.Config, notifier *observer.Notifier) *EventGenerator {
	return &EventGenerator{
		config:      config,
		competitors: make(map[int]*competitor),
		notifier:    notifier,
	}
}

// Notify - локально обновляет состояние участников и генерируя новые события.
func (g *EventGenerator) Notify(event events.EventData) {
	// Проверка всех участников на превышение допустимого времени старта
	for id, comp := range g.competitors {
		if !comp.started &&
			comp.startTime.Add(g.config.StartDelta).Before(event.Time) &&
			!comp.disqualified {
			// Автоматическая дисквалификация при нарушении временных рамок
			comp.disqualified = true

			newEvent := event
			newEvent.ID = events.EventDisqualified
			newEvent.Params = ""
			newEvent.CompetitorID = id

			g.notifier.NotifyAll(newEvent)
		}
	}

	comp, exists := g.competitors[event.CompetitorID]

	switch event.ID {
	case events.EventRegistered:
		// Регистрация нового участника
		if !exists {
			// Инициализация с "максимально поздним" временем старта,
			// чтобы не дисквалифировали системой до получения реального времени старта
			maxTime := time.Date(0, 1, 1, 23, 59, 59, 0, time.UTC)
			g.competitors[event.CompetitorID] = &competitor{
				endedLaps: 0,
				startTime: maxTime.Add(-g.config.StartDelta),
			}
		} else {
			fmt.Printf("Competitor(%d) is registered\n", event.CompetitorID)
			return
		}

	case events.EventStartTimeSet:
		// Установка реального времени старта
		if !exists {
			fmt.Printf("Competitor(%d) is not registered\n", event.CompetitorID)
			return
		}

		startTime, err := utils.ParseStrTimeToTime(event.Params)
		if err != nil {
			fmt.Printf("%s", err)
		}

		comp.startTime = startTime // Обновление времени старта

	case events.EventStarted:
		// Обработка фактического старта участника
		if !exists {
			fmt.Printf("Competitor(%d) is not registered\n", event.CompetitorID)
			return
		}

		// Проверка на досрочный старт
		if event.Time.Before(comp.startTime) {
			comp.disqualified = true

			newEvent := event
			newEvent.ID = events.EventDisqualified
			newEvent.Params = ""

			g.notifier.NotifyAll(newEvent)
		}

		comp.started = true // Флаг начала гонки

	case events.EventLapEnded:
		// Обработка завершения круга
		if !exists {
			fmt.Printf("Competitor(%d) is not registered\n", event.CompetitorID)
			return
		}

		comp.endedLaps += 1

		// Проверка что круг последний
		if comp.endedLaps == g.config.Laps {
			newEvent := event
			newEvent.ID = events.EventFinished
			g.notifier.NotifyAll(newEvent)
		}
	}
}
