package observer

import "biathlonCompetitions/internal/events"

// Observer определяет интерфейс для подписчиков на события в паттерне "Наблюдатель".
// Любой тип, реализующий этот интерфейс, может получать уведомления о событиях гонки.
type Observer interface {
	Notify(event events.EventData)
}
