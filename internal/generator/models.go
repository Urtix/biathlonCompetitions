package generator

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/observer"
	"time"
)

// competitor содержит состояние участника гонки в реальном времени
type competitor struct {
	endedLaps    int       // Количество успешно завершенных кругов
	startTime    time.Time // Персональное время старта участника
	started      bool      // Флаг начала гонки (true = участник стартовал)
	disqualified bool      // Статус дисквалификации (true = дисквалифицирован)
}

// EventGenerator - управление событиями гонки:
// 1. Отслеживает состояние всех участников
// 2. Автоматически генерирует исходящие события(когда это требуется)
type EventGenerator struct {
	config      config.Config       // Конфигурация параметров гонки
	competitors map[int]*competitor // Текущие участники (ID -> состояние)
	notifier    *observer.Notifier  // Система рассылки событий (паттерн Observer)
}
