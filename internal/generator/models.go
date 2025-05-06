package generator

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/observer"
	"time"
)

type competitor struct {
	endedLaps    int
	startTime    time.Time
	started      bool
	disqualified bool
}

type EventGenerator struct {
	config      config.Config
	competitors map[int]*competitor
	notifier    *observer.Notifier
}
