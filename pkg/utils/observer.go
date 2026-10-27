package utils

// Observer is a simple implementation of the observer pattern. It allows you to subscribe to events and notify subscribers when an event occurs.
type Observer[T any] struct {
	observers []func(T)
}

func (o *Observer[T]) Subscribe(observer func(T)) {
	o.observers = append(o.observers, observer)
}

func (o *Observer[T]) Notify(data T) {
	for _, observer := range o.observers {
		observer(data)
	}
}

func (o *Observer[T]) Clear() {
	o.observers = nil
}
