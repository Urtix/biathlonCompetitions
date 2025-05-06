package handler

import (
	"biathlonCompetitions/internal/config"
	"time"
)

type lapInfo struct {
	duration time.Duration
	speed    float64
}

type mainLaps struct {
	totalDuration string
	startTime     time.Time // время старта для текущего круга, а не общее
	info          []lapInfo
}

type penaltyLaps struct {
	currentLap int
	startTime  time.Time
	info       lapInfo
}

type shoot struct {
	currentTargets [5]int
	hitCount       int
	shotCount      int
}

type competitor struct {
	id           int
	disqualified bool
	mainLaps     mainLaps
	penaltyLaps  penaltyLaps
	shoot        shoot
}

// Tracker реализует Observer
type Tracker struct {
	config      config.Config
	competitors map[int]*competitor
}
