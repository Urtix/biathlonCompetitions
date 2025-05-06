package observer

import (
	"biathlonCompetitions/internal/events"
)

// Notifier реализует паттерн Наблюдатель (Observer).
// Обеспечивает механизм подписки и рассылки событий между компонентами системы.
type Notifier struct {
	observers []Observer // Список зарегистрированных наблюдателей
}

// Register добавляет нового наблюдателя в список подписчиков.
// Принимает объект, реализующий интерфейс Observer.
func (n *Notifier) Register(observer Observer) {
	n.observers = append(n.observers, observer)
}

// NotifyAll рассылает событие всем зарегистрированным наблюдателям.
// Каждый наблюдатель получает полную копию события для обработки.
func (n *Notifier) NotifyAll(event events.EventData) {
	for _, observer := range n.observers {
		observer.Notify(event)
	}
}
