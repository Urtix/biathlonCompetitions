package handler

import (
	"biathlonCompetitions/internal/config"
	"time"
)

// lapInfo содержит данные о пройденном круге
type lapInfo struct {
	duration time.Duration // Затраченное время
	speed    float64       // Средняя скорость (м/с)
}

// mainLaps хранит информацию об основных кругах гонки
type mainLaps struct {
	totalDuration string    // Общее время гонки (строковое представление)
	startTime     time.Time // Время старта текущего круга
	info          []lapInfo // Информация по каждому пройденному кругу
}

// penaltyLaps - штрафные круги участника
type penaltyLaps struct {
	currentLap int       // Номер текущего штрафного круаг
	startTime  time.Time // Время начала штрафного круга
	info       lapInfo   // Общая информаци о штрафных кругах
}

// shoot отслеживает результаты стрельбы
type shoot struct {
	currentTargets [5]bool // Текущие мишени (true - поражена, false - нет)
	hitCount       int     // Общее количество попаданий
	shotCount      int     // Всего произведено выстрелов
}

// competitor содержит полное состояние участника гонки
type competitor struct {
	id           int         // Уникальный идентификатор
	disqualified bool        // Флаг дисквалификации
	mainLaps     mainLaps    // Основные круги
	penaltyLaps  penaltyLaps // Штрафные круги
	shoot        shoot       // Результаты стрельбы
}

// Handler реализует паттерн Observer для отслеживания событий гонки
// Хранит состояние всех участников и конфигурацию соревнований
type Handler struct {
	config      config.Config       // Параметры гонки
	competitors map[int]*competitor // Участники (ID → состояние)
}
