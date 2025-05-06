package config

import (
	"biathlonCompetitions/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// LoadConfig загружает и валидирует конфигурацию приложения из JSON-файла
func LoadConfig(filename string) (Config, error) {
	var jsonConf jsonConfig

	// Чтение всего содержимого файла в память
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("%w: %s", ErrFileNotFound, filename)
		}
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	// Парсинг JSON данных во временную структуру
	if err = json.Unmarshal(data, &jsonConf); err != nil {
		return Config{}, fmt.Errorf("config parsing failed: %w", err)
	}

	// Валидация числовых параметров конфигурации
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

	// Преобразование time.Time.string()  в тип time.Time
	startTime, err := utils.ParseStrTimeToTime(jsonConf.Start)
	if err != nil {
		return Config{}, fmt.Errorf("%w: start time: %v", ErrInvalidTimeFormat, err)
	}

	// Преобразование time.Time.string() в time.Duration
	deltaTime, err := utils.ParseStrTimeToDuration(jsonConf.StartDelta)
	if err != nil {
		return Config{}, fmt.Errorf("%w: start delta time: %v", ErrInvalidTimeFormat, err)
	}

	// Сборка финального объекта конфигурации
	return Config{
		Laps:        jsonConf.Laps,
		LapLen:      jsonConf.LapLen,
		PenaltyLen:  jsonConf.PenaltyLen,
		FiringLines: jsonConf.FiringLines,
		Start:       startTime,
		StartDelta:  deltaTime,
	}, nil
}
