package observer

import (
	"biathlonCompetitions/internal/events"
)

type Notifier struct {
	observers []Observer
}

func (n *Notifier) Register(observer Observer) {
	n.observers = append(n.observers, observer)
}

func (n *Notifier) NotifyAll(event events.EventData) {
	for _, observer := range n.observers {
		observer.Notify(event)
	}
}
