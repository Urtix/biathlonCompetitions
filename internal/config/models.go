package config

import (
	"errors"
	"time"
)

type jsonConfig struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type Config struct {
	Laps        int
	LapLen      int
	PenaltyLen  int
	FiringLines int
	Start       time.Time
	StartDelta  time.Duration
}

var (
	ErrFileNotFound      = errors.New("config file not found")
	ErrInvalidTimeFormat = errors.New("invalid time format")
	ErrInvalidValue      = errors.New("invalid value")
)
