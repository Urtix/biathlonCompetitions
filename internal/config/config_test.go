package config

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestConfig_LoadConfig(t *testing.T) {
	createTestFile := func(content string) string {
		tmpFile, err := os.CreateTemp("", "test-config-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer tmpFile.Close()

		if _, err = tmpFile.WriteString(content); err != nil {
			t.Fatal(err)
		}
		return tmpFile.Name()
	}

	tests := []struct {
		name       string
		wantErr    error
		wantConfig Config
		setup      func() string // Возвращает имя файла
		cleanup    func(string)
	}{
		{
			name: "Valid config",
			wantConfig: Config{
				Laps:        3,
				LapLen:      1000,
				PenaltyLen:  200,
				FiringLines: 5,
				Start:       parseTime("09:00:00.000"),
				StartDelta:  5 * time.Minute,
			},
			setup: func() string {
				return createTestFile(`{
                    "laps": 3,
                    "lapLen": 1000,
                    "penaltyLen": 200,
                    "firingLines": 5,
                    "start": "09:00:00.000",
                    "startDelta": "00:05:00"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
		{
			name:    "File not found",
			wantErr: ErrFileNotFound,
			setup:   func() string { return "non-existent-file.json" },
			cleanup: func(f string) {},
		},
		{
			name:    "Invalid laps",
			wantErr: ErrInvalidValue,
			setup: func() string {
				return createTestFile(`{
                    "laps": 0,
                    "lapLen": 1000,
                    "penaltyLen": 200,
                    "firingLines": 5,
                    "start": "09:00:00.000",
                    "startDelta": "00:05:00"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
		{
			name:    "Invalid lap len",
			wantErr: ErrInvalidValue,
			setup: func() string {
				return createTestFile(`{
                    "laps": 3,
                    "lapLen": 0,
                    "penaltyLen": 200,
                    "firingLines": 5,
                    "start": "09:00:00.000",
                    "startDelta": "00:05:00"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
		{
			name:    "Invalid penalty len",
			wantErr: ErrInvalidValue,
			setup: func() string {
				return createTestFile(`{
                    "laps": 3,
                    "lapLen": 1000,
                    "penaltyLen": 0,
                    "firingLines": 5,
                    "start": "09:00:00.000",
                    "startDelta": "00:05:00"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
		{
			name:    "Invalid firing lines",
			wantErr: ErrInvalidValue,
			setup: func() string {
				return createTestFile(`{
                    "laps": 3,
                    "lapLen": 1000,
                    "penaltyLen": 200,
                    "firingLines": 0,
                    "start": "09:00:00.000",
                    "startDelta": "00:05:00"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
		{
			name:    "Invalid start time format",
			wantErr: ErrInvalidTimeFormat,
			setup: func() string {
				return createTestFile(`{
                    "laps": 3,
                    "lapLen": 1000,
                    "penaltyLen": 200,
                    "firingLines": 5,
                    "start": "time",
                    "startDelta": "00:05:00"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
		{
			name:    "Invalid start delta format",
			wantErr: ErrInvalidTimeFormat,
			setup: func() string {
				return createTestFile(`{
                    "laps": 3,
                    "lapLen": 1000,
                    "penaltyLen": 200,
                    "firingLines": 5,
                    "start": "09:00:00.000",
                    "startDelta": "delta"
                }`)
			},
			cleanup: func(f string) { os.Remove(f) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setup()
			defer tt.cleanup(filename)

			cfg, err := LoadConfig(filename)

			// Проверка ошибок
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Проверка значений конфига
			if cfg.Laps != tt.wantConfig.Laps {
				t.Errorf("Laps: got %d, want %d", cfg.Laps, tt.wantConfig.Laps)
			}

			if cfg.LapLen != tt.wantConfig.LapLen {
				t.Errorf("LapLen: got %d, want %d", cfg.LapLen, tt.wantConfig.LapLen)
			}

			if !cfg.Start.Equal(tt.wantConfig.Start) {
				t.Errorf("Start: got %v, want %v", cfg.Start, tt.wantConfig.Start)
			}

			if cfg.StartDelta != tt.wantConfig.StartDelta {
				t.Errorf("StartDelta: got %v, want %v", cfg.StartDelta, tt.wantConfig.StartDelta)
			}
		})
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("15:04:05.000", s)
	return t
}
