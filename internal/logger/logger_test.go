package logger

import (
	"biathlonCompetitions/internal/events"
	"testing"
)

func TestGenerateMessage(t *testing.T) {
	tests := []struct {
		name     string
		event    events.EventData
		expected string
	}{
		{
			name: "EventRegistered",
			event: events.EventData{
				ID:           events.EventRegistered,
				CompetitorID: 1,
			},
			expected: "The competitor(1) registered",
		},
		{
			name: "EventStartTimeSet",
			event: events.EventData{
				ID:           events.EventStartTimeSet,
				CompetitorID: 1,
				Params:       "12:30:00.000",
			},
			expected: "The start time for the competitor(1) was set by a draw to 12:30:00.000",
		},
		{
			name: "EventOnStartLine",
			event: events.EventData{
				ID:           events.EventOnStartLine,
				CompetitorID: 1,
			},
			expected: "The competitor(1) is on the start line",
		},
		{
			name: "EventStarted",
			event: events.EventData{
				ID:           events.EventStarted,
				CompetitorID: 1,
			},
			expected: "The competitor(1) has started",
		},
		{
			name: "EventOnFiringRange",
			event: events.EventData{
				ID:           events.EventOnFiringRange,
				CompetitorID: 1,
				Params:       "1",
			},
			expected: "The competitor(1) is on the firing range(1)",
		},
		{
			name: "EventTargetHit",
			event: events.EventData{
				ID:           events.EventTargetHit,
				CompetitorID: 1,
				Params:       "1",
			},
			expected: "The target(1) has been hit by competitor(1)",
		},
		{
			name: "EventLeftFiringRange",
			event: events.EventData{
				ID:           events.EventLeftFiringRange,
				CompetitorID: 1,
			},
			expected: "The competitor(1) left the firing range",
		},
		{
			name: "EventEnteredPenalty",
			event: events.EventData{
				ID:           events.EventEnteredPenalty,
				CompetitorID: 1,
			},
			expected: "The competitor(1) entered the penalty laps",
		},
		{
			name: "EventLeftPenalty",
			event: events.EventData{
				ID:           events.EventLeftPenalty,
				CompetitorID: 1,
			},
			expected: "The competitor(1) left the penalty laps",
		},
		{
			name: "EventLapEnded",
			event: events.EventData{
				ID:           events.EventLapEnded,
				CompetitorID: 1,
			},
			expected: "The competitor(1) ended the main lap",
		},
		{
			name: "EventCannotContinue",
			event: events.EventData{
				ID:           events.EventCannotContinue,
				CompetitorID: 1,
				Params:       "Lost in the forest",
			},
			expected: "The competitor(1) can't continue: Lost in the forest",
		},
		{
			name: "EventDisqualified",
			event: events.EventData{
				ID:           events.EventDisqualified,
				CompetitorID: 1,
			},
			expected: "The competitor(1) is disqualified",
		},
		{
			name: "EventFinished",
			event: events.EventData{
				ID:           events.EventFinished,
				CompetitorID: 1,
			},
			expected: "The competitor(1) has finished",
		},
		{
			name: "UnknownEvent",
			event: events.EventData{
				ID:           -1,
				CompetitorID: 1,
			},
			expected: "Unknown event(-1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateMessage(tt.event)
			if result != tt.expected {
				t.Errorf("\nTest: %s\nExpected: %q\nActual:   %q",
					tt.name, tt.expected, result)
			}
		})
	}
}
