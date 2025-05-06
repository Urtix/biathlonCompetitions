package events

import "time"

const (
	EventRegistered      = 1  // Участник зарегистрирован в системе
	EventStartTimeSet    = 2  // Установлено время старта для участника
	EventOnStartLine     = 3  // Участник занял позицию на стартовой линии
	EventStarted         = 4  // Участник начал гонку
	EventOnFiringRange   = 5  // Участник прибыл на стрелковый рубеж
	EventTargetHit       = 6  // Успешное попадание в мишень
	EventLeftFiringRange = 7  // Участник покинул стрелковый рубеж
	EventEnteredPenalty  = 8  // Участник начал штрафной круг
	EventLeftPenalty     = 9  // Участник завершил штрафной круг
	EventLapEnded        = 10 // Участник завершил основной круг
	EventCannotContinue  = 11 // Участник не может продолжать гонку
	EventDisqualified    = 32 // Участник дисквалифицирован
	EventFinished        = 33 // Участник успешно завершил гонку
)

// EventData содержит полную информацию о текущем событии
type EventData struct {
	ID           int       // Идентификатор события
	CompetitorID int       // Идентификатор участника
	Params       string    // Дополнительные параметры
	Time         time.Time // Время события
}
