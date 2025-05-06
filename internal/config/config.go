package config

import (
	"biathlonCompetitions/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func LoadConfig(filename string) (Config, error) {
	var jsonConf jsonConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("%w: %s", ErrFileNotFound, filename)
		}
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	if err = json.Unmarshal(data, &jsonConf); err != nil {
		return Config{}, fmt.Errorf("config parsing failed: %w", err)
	}

	if jsonConf.Laps < 1 {
		return Config{}, fmt.Errorf("%w: laps must be >= 1", ErrInvalidValue)
	}

	if jsonConf.LapLen < 1 {
		return Config{}, fmt.Errorf("%w: lapLen must be >= 1", ErrInvalidValue)
	}

	if jsonConf.PenaltyLen < 1 {
		return Config{}, fmt.Errorf("%w: penaltyLen must be >= 1", ErrInvalidValue)
	}

	if jsonConf.FiringLines < 1 {
		return Config{}, fmt.Errorf("%w: firingLines must be >= 1", ErrInvalidValue)
	}

	startTime, err := utils.ParseStrTimeToTime(jsonConf.Start)
	if err != nil {
		return Config{}, fmt.Errorf("%w: start time: %v", ErrInvalidTimeFormat, err)
	}

	deltaTime, err := utils.ParseStrTimeToDuration(jsonConf.StartDelta)
	if err != nil {
		return Config{}, fmt.Errorf("%w: start delta time: %v", ErrInvalidTimeFormat, err)
	}

	return Config{
		Laps:        jsonConf.Laps,
		LapLen:      jsonConf.LapLen,
		PenaltyLen:  jsonConf.PenaltyLen,
		FiringLines: jsonConf.FiringLines,
		Start:       startTime,
		StartDelta:  deltaTime,
	}, nil
}
