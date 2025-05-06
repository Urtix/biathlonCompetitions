package config

import (
	"errors"
	"time"
)

// jsonConfig - временная структура для парсинга JSON-файла конфигурации.
// Используется только для десериализации перед валидацией и преобразованием в Config.
type jsonConfig struct {
	Laps        int    `json:"laps"`        // Общее количество кругов в гонке
	LapLen      int    `json:"lapLen"`      // Длина одного круга в метрах
	PenaltyLen  int    `json:"penaltyLen"`  // Длина штрафного круга в метрах
	FiringLines int    `json:"firingLines"` // Количество стрелковых рубежей
	Start       string `json:"start"`       // Планируемое время старта (формат "15:04")
	StartDelta  string `json:"startDelta"`  // Максимальное смещение старта (формат "1h30m")
}

// Config - валидированная конфигурация со всеми параметрами.
// Содержит преобразованные к итоговым типам значения из jsonConfig.
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
