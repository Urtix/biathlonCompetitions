package handler

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/events"
	"testing"
	"time"
)

func TestHandler_Notify(t *testing.T) {
	cfg := config.Config{
		Laps:        2,
		LapLen:      600,
		PenaltyLen:  200,
		FiringLines: 2,
		Start:       time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		StartDelta:  1*time.Minute + 30*time.Second,
	}

	tests := []struct {
		name     string
		events   []events.EventData
		validate func(*testing.T, *Handler)
	}{
		{
			name: "Shooting first firing lines",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 1},
				{ID: events.EventOnFiringRange, CompetitorID: 1, Params: "1"},
				{ID: events.EventTargetHit, CompetitorID: 1, Params: "1"},
				{ID: events.EventLeftFiringRange, CompetitorID: 1},
			},
			validate: func(t *testing.T, tr *Handler) {
				comp := tr.competitors[1]
				if comp.shoot.hitCount != 1 {
					t.Errorf("Expected 1 hit, got %d", comp.shoot.hitCount)
				}
				if comp.shoot.shotCount != 5 {
					t.Errorf("Expected 5 shots, got %d", comp.shoot.shotCount)
				}
			},
		},
		{
			name: "Shooting second firing lines",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 1},
				{ID: events.EventOnFiringRange, CompetitorID: 1, Params: "1"},
				{ID: events.EventTargetHit, CompetitorID: 1, Params: "1"},
				{ID: events.EventLeftFiringRange, CompetitorID: 1},
				{ID: events.EventOnFiringRange, CompetitorID: 1, Params: "2"},
				{ID: events.EventLeftFiringRange, CompetitorID: 1},
			},
			validate: func(t *testing.T, tr *Handler) {
				comp := tr.competitors[1]
				if comp.shoot.hitCount != 1 {
					t.Errorf("Expected 1 hit, got %d", comp.shoot.hitCount)
				}
				if comp.shoot.shotCount != 10 {
					t.Errorf("Expected 10 shots, got %d", comp.shoot.shotCount)
				}
			},
		},
		{
			name: "Penalty",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 3},
				{ID: events.EventEnteredPenalty, CompetitorID: 3, Time: parseTime("10:00:00.000")},
				{ID: events.EventLeftPenalty, CompetitorID: 3, Time: parseTime("10:00:20.000")},
			},
			validate: func(t *testing.T, tr *Handler) {
				comp := tr.competitors[3]
				if comp.penaltyLaps.info.duration != 20*time.Second {
					t.Errorf("Expected 30s penalty, got %v", comp.penaltyLaps.info.duration)
				}
				if comp.penaltyLaps.info.speed != 10 {
					t.Errorf("Expected 10 m/s penalty, got %v", comp.penaltyLaps.info.speed)
				}
			},
		},
		{
			name: "Cannot continue",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 1},
				{ID: events.EventStartTimeSet, CompetitorID: 1, Params: "10:00:00.000"},
				{ID: events.EventStarted, CompetitorID: 1, Time: parseTime("10:00:01.000")},
				{ID: events.EventCannotContinue, CompetitorID: 1},
			},
			validate: func(t *testing.T, tr *Handler) {
				comp := tr.competitors[1]

				if comp.mainLaps.totalDuration != "NotFinished" {
					t.Errorf("Expected NotFinished, got %v", comp.mainLaps.totalDuration)
				}
			},
		},
		{
			name: "Full race cycle",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 2},
				{ID: events.EventStartTimeSet, CompetitorID: 2, Params: "10:00:00.000"},
				{ID: events.EventStarted, CompetitorID: 2, Time: parseTime("10:00:01.000")},
				{ID: events.EventOnFiringRange, CompetitorID: 2, Time: parseTime("10:00:10.000")},
				{ID: events.EventTargetHit, CompetitorID: 2, Time: parseTime("10:00:20.000"), Params: "1"},
				{ID: events.EventTargetHit, CompetitorID: 2, Time: parseTime("10:00:20.000"), Params: "2"},
				{ID: events.EventTargetHit, CompetitorID: 2, Time: parseTime("10:00:20.000"), Params: "3"},
				{ID: events.EventTargetHit, CompetitorID: 2, Time: parseTime("10:00:20.000"), Params: "4"},
				{ID: events.EventTargetHit, CompetitorID: 2, Time: parseTime("10:00:20.000"), Params: "5"},
				{ID: events.EventLeftFiringRange, CompetitorID: 2, Time: parseTime("10:00:30.000")},
				{ID: events.EventLapEnded, CompetitorID: 2, Time: parseTime("10:01:00.000")},
				{ID: events.EventOnFiringRange, CompetitorID: 2, Time: parseTime("10:01:10.000")},
				{ID: events.EventLeftFiringRange, CompetitorID: 2, Time: parseTime("10:01:20.000")},
				{ID: events.EventEnteredPenalty, CompetitorID: 2, Time: parseTime("10:01:40.000")},
				{ID: events.EventLeftPenalty, CompetitorID: 2, Time: parseTime("10:01:50.000")},
				{ID: events.EventLapEnded, CompetitorID: 2, Time: parseTime("10:02:00.000")},
				{ID: events.EventFinished, CompetitorID: 2},
			},
			validate: func(t *testing.T, tr *Handler) {
				comp := tr.competitors[2]

				if len(comp.mainLaps.info) != 2 {
					t.Errorf("Expected 2 laps, got %d", len(comp.mainLaps.info))
				}

				if comp.mainLaps.totalDuration != "2m0s" {
					t.Errorf("Expected 5 total time, got %s", comp.mainLaps.totalDuration)
				}

				if comp.mainLaps.info[0].duration != 1*time.Minute {
					t.Errorf("Expected 1m0s, got %d", comp.mainLaps.info[0].duration)
				}
				if comp.mainLaps.info[0].speed != 10. {
					t.Errorf("Expected 10 m/s, got %f m/s", comp.mainLaps.info[0].speed)
				}
				if comp.mainLaps.info[1].duration != 1*time.Minute {
					t.Errorf("Expected 1m0s, got %d", comp.mainLaps.info[0].duration)
				}
				if comp.mainLaps.info[1].speed != 10. {
					t.Errorf("Expected 10 m/s, got %f m/s", comp.mainLaps.info[1].speed)
				}

				if comp.penaltyLaps.info.duration != 10*time.Second {
					t.Errorf("Expected 10s, got %d", comp.penaltyLaps.info.duration)
				}
				if comp.penaltyLaps.info.speed != 20. {
					t.Errorf("Expected 20 m/s, got %f m/s", comp.penaltyLaps.info.speed)
				}

				if comp.shoot.hitCount != 5 {
					t.Errorf("Expected 5 hits, got %d", comp.shoot.hitCount)
				}
				if comp.shoot.shotCount != 10 {
					t.Errorf("Expected 5 shots, got %d", comp.shoot.shotCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewHandler(cfg)
			for _, event := range tt.events {
				tr.Notify(event)
			}
			tt.validate(t, tr)
		})
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("15:04:05.000", s)
	return t
}
