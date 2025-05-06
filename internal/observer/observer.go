package observer

import "biathlonCompetitions/internal/events"

type Observer interface {
	Notify(event events.EventData)
}
