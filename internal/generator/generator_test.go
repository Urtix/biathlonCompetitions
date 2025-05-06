package generator

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/events"
	"biathlonCompetitions/internal/observer"
	"testing"
	"time"
)

func TestGenerator_Notify(t *testing.T) {
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
		validate func(*testing.T, *EventGenerator)
	}{
		{
			name: "False start",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 1},
				{ID: events.EventStartTimeSet, CompetitorID: 1, Time: parseTime("10:20:00.000"), Params: "10:30:00.000"},
				{ID: events.EventStarted, CompetitorID: 1, Time: parseTime("10:29:59.999")},
			},
			validate: func(t *testing.T, tr *EventGenerator) {
				comp := tr.competitors[1]
				if !comp.disqualified {
					t.Errorf("Expected true, got %v", comp.disqualified)
				}
			},
		},
		{
			name: "Not disqualified",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 1},
				{ID: events.EventStartTimeSet, CompetitorID: 1, Time: parseTime("10:20:00.000"), Params: "10:30:00.000"},
				{ID: events.EventStarted, CompetitorID: 1, Time: parseTime("10:30:00.000")},
			},
			validate: func(t *testing.T, tr *EventGenerator) {
				comp := tr.competitors[1]
				if comp.disqualified {
					t.Errorf("Expected false, got %v", comp.disqualified)
				}
			},
		},
		{
			name: "Not start in delta start",
			events: []events.EventData{
				{ID: events.EventRegistered, CompetitorID: 1, Time: parseTime("09:00:00.000")},
				{ID: events.EventStartTimeSet, CompetitorID: 1, Time: parseTime("09:20:00.000"), Params: "10:00:00.000"},
				{ID: events.EventRegistered, CompetitorID: 2, Time: parseTime("10:01:30.001")},
			},
			validate: func(t *testing.T, tr *EventGenerator) {
				comp := tr.competitors[1]
				if !comp.disqualified {
					t.Errorf("Expected true, got %v", comp.disqualified)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := &observer.Notifier{}
			eg := NewEventGenerator(cfg, notifier)
			for _, event := range tt.events {
				eg.Notify(event)
			}
			tt.validate(t, eg)
		})
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("15:04:05.000", s)
	return t
}
