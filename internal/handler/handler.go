package handler

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/events"
	"biathlonCompetitions/internal/utils"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func NewTracker(config config.Config) *Tracker {
	return &Tracker{
		config:      config,
		competitors: make(map[int]*competitor),
	}
}

func (t *Tracker) Notify(event events.EventData) {
	comp, exists := t.competitors[event.CompetitorID]

	switch event.ID {
	case events.EventRegistered:
		if !exists {
			comp = &competitor{
				id: event.CompetitorID,
				mainLaps: mainLaps{
					totalDuration: "NotStarted",
				},
			}
			t.competitors[event.CompetitorID] = comp
		} else {
			fmt.Printf("Участник (%d) уже зарегестрирован\n", event.CompetitorID)
		}

	case events.EventStartTimeSet:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		startLapTime, err := time.Parse("15:04:05.000", event.Params)
		if err != nil {
			fmt.Printf("%s: invalid duration", err)
			return
		}

		comp.mainLaps.startTime = startLapTime

	case events.EventOnStartLine:

	case events.EventStarted:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		if event.Time.Before(comp.mainLaps.startTime) || event.Time.After(comp.mainLaps.startTime.Add(t.config.StartDelta)) {
			comp.disqualified = true
		}

	case events.EventOnFiringRange:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		comp.shoot.shotCount += 5

	case events.EventTargetHit:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		idTargetHit, err := strconv.Atoi(event.Params)
		if err != nil {
			fmt.Printf("%s: invalid target id", err)
			return
		}

		if idTargetHit < 1 || idTargetHit > len(comp.shoot.currentTargets) {
			fmt.Printf("invalid target id: %d", idTargetHit)
			return
		}

		comp.shoot.currentTargets[idTargetHit-1] = 1

	case events.EventLeftFiringRange:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		comp.shoot.hitCount += utils.SumIntArray(comp.shoot.currentTargets[:])
		comp.shoot.currentTargets = [5]int{}

	case events.EventEnteredPenalty:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		comp.penaltyLaps.currentLap += 1
		comp.penaltyLaps.startTime = event.Time

	case events.EventLeftPenalty:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		comp.penaltyLaps.info.duration += event.Time.Sub(comp.penaltyLaps.startTime)
		totalPenaltyLen := t.config.PenaltyLen * comp.penaltyLaps.currentLap
		comp.penaltyLaps.info.speed = float64(totalPenaltyLen) / comp.penaltyLaps.info.duration.Seconds()

	case events.EventLapEnded:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		lapTime := event.Time.Sub(comp.mainLaps.startTime)
		newLapInfo := lapInfo{
			duration: lapTime,
			speed:    float64(t.config.LapLen) / lapTime.Seconds(),
		}

		comp.mainLaps.info = append(comp.mainLaps.info, newLapInfo)
		comp.mainLaps.startTime = event.Time

	case events.EventCannotContinue:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		comp.mainLaps.totalDuration = "NotFinished"
		comp.disqualified = true

	case events.EventDisqualified:
		comp.mainLaps.totalDuration = "NotStarted"
		comp.disqualified = true

	case events.EventFinished:
		if !exists {
			fmt.Printf("Участника (%d) не существует\n", event.CompetitorID)
			return
		}

		if comp.disqualified {
			return
		}

		var totalTime time.Duration
		for _, lap := range comp.mainLaps.info {
			totalTime += lap.duration
		}

		comp.mainLaps.totalDuration = totalTime.String()
	}
}

func (t *Tracker) FinalReport() {
	competitors := make([]*competitor, 0, len(t.competitors))
	for _, comp := range t.competitors {
		competitors = append(competitors, comp)
	}

	sort.Slice(competitors, func(i, j int) bool {
		return competitors[i].mainLaps.totalDuration < competitors[j].mainLaps.totalDuration
	})

	for _, comp := range competitors {
		totalStrTime := utils.ParseStrDurationToStrTime(comp.mainLaps.totalDuration)
		penaltyTimeStr := utils.ParseDurationToStrTime(comp.penaltyLaps.info.duration)

		fmt.Printf("[%s] %d", totalStrTime, comp.id)
		fmt.Printf(" [")

		for i := 0; i < t.config.Laps; i++ {
			if i < len(comp.mainLaps.info) {
				lapTimeStr := utils.ParseDurationToStrTime(comp.mainLaps.info[i].duration)
				fmt.Printf("{%s, %0.3f}", lapTimeStr, comp.mainLaps.info[i].speed)
			} else {
				fmt.Printf("{ , }")
			}

			if i < t.config.Laps-1 {
				fmt.Printf(", ")
			}
		}

		fmt.Printf("] {%s, %0.3f} %d/%d\n",
			penaltyTimeStr,
			comp.penaltyLaps.info.speed,
			comp.shoot.hitCount,
			comp.shoot.shotCount,
		)
	}
}
