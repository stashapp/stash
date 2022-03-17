package manager

import (
	"context"
	"sync"
)

type subscriptionManager struct {
	subscriptions []chan bool
	mutex         sync.Mutex
}

func (m *subscriptionManager) subscribe(ctx context.Context) <-chan bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	c := make(chan bool, 10)
	m.subscriptions = append(m.subscriptions, c)

	go func() {
		<-ctx.Done()
		m.mutex.Lock()
		defer m.mutex.Unlock()
		close(c)

		for i, s := range m.subscriptions {
			if s == c {
				m.subscriptions = append(m.subscriptions[:i], m.subscriptions[i+1:]...)
				break
			}
		}
	}()

	return c
}

func (m *subscriptionManager) notify() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, s := range m.subscriptions {
		s <- true
	}
}
