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

func (g *EventGenerator) Notify(event events.EventData) {
	for id, comp := range g.competitors {
		if !comp.started && comp.startTime.Add(g.config.StartDelta).Before(event.Time) && !comp.disqualified {
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
		if !exists {
			maxTime := time.Date(0, 1, 1, 23, 59, 59, 0, time.UTC)
			g.competitors[event.CompetitorID] =
				&competitor{
					endedLaps: 0,
					startTime: maxTime.Add(-g.config.StartDelta),
				}
		} else {
			fmt.Printf("Участник (%d) уже зарегестрирован\n", event.CompetitorID)
			return
		}

	case events.EventStartTimeSet:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		startTime, err := utils.ParseStrTimeToTime(event.Params)
		if err != nil {
			fmt.Printf("%s", err)
		}

		comp.startTime = startTime

	case events.EventStarted:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		if event.Time.Before(comp.startTime) {
			comp.disqualified = true

			newEvent := event
			newEvent.ID = events.EventDisqualified
			newEvent.Params = ""

			g.notifier.NotifyAll(newEvent)
		}

		comp.started = true

	case events.EventLapEnded:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		comp.endedLaps += 1

		if comp.endedLaps == g.config.Laps {
			newEvent := event
			newEvent.ID = events.EventFinished
			g.notifier.NotifyAll(newEvent)
		}
	}
}
