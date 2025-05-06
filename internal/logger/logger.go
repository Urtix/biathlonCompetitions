package logger

import (
	"biathlonCompetitions/internal/events"
	"fmt"
)

func NewLogger() *Logger {
	return &Logger{}
}

func (j *Logger) Notify(event events.EventData) {
	message := generateMessage(event)
	fmt.Printf("[%s] %s\n", event.Time.Format("15:04:05.000"), message)
}

func generateMessage(event events.EventData) string {
	var message string

	switch event.ID {
	case events.EventRegistered:
		message = fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)
	case events.EventStartTimeSet:
		message = fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", event.CompetitorID, event.Params)
	case events.EventOnStartLine:
		message = fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
	case events.EventStarted:
		message = fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
	case events.EventOnFiringRange:
		message = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.Params)
	case events.EventTargetHit:
		message = fmt.Sprintf("The target(%s) has been hit by competitor(%d)", event.Params, event.CompetitorID)
	case events.EventLeftFiringRange:
		message = fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
	case events.EventEnteredPenalty:
		message = fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
	case events.EventLeftPenalty:
		message = fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
	case events.EventLapEnded:
		message = fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
	case events.EventCannotContinue:
		message = fmt.Sprintf("The competitor(%d) can't continue: %s", event.CompetitorID, event.Params)
	case events.EventDisqualified:
		message = fmt.Sprintf("The competitor(%d) is disqualified", event.CompetitorID)
	case events.EventFinished:
		message = fmt.Sprintf("The competitor(%d) has finished", event.CompetitorID)
	default:
		message = fmt.Sprintf("Unknown event(%d)", event.ID)
	}

	return message
}
