package events

import "time"

const (
	EventRegistered      = 1
	EventStartTimeSet    = 2
	EventOnStartLine     = 3
	EventStarted         = 4
	EventOnFiringRange   = 5
	EventTargetHit       = 6
	EventLeftFiringRange = 7
	EventEnteredPenalty  = 8
	EventLeftPenalty     = 9
	EventLapEnded        = 10
	EventCannotContinue  = 11
	EventDisqualified    = 32
	EventFinished        = 33
)

type EventData struct {
	ID           int
	CompetitorID int
	Params       string
	Time         time.Time
}
