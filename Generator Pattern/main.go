package main

import (
	"fmt"
	"sync"
	"time"
)

type (
	Event struct {
		data int
	}

	Observer interface {
		NotifyCallback(Event)
	}

	Subject interface {
		AddListener(observer Observer)
		RemoveListener(observer Observer)
		Notify(event Event)
	}

	eventObserver struct {
		id int
		time time.Time
	}

	eventSubject struct {
		observers sync.Map
	}
)

func (e *eventObserver) NotifyCallback(event Event) {
	fmt.Printf("received %d after %v\n", event.data, time.Since(e.time))
}

func (s *eventSubject) AddListener(obs Observer) {
	s.observers.Store(obs, struct{}{})
}

func (s *eventSubject) RemoveListener(obs Observer) {
	s.observers.Delete(obs)
}

func (s *eventSubject) Nofity(event Event) {
	s.observers.Range(func(key interface{}, value interface{}) bool {
		if key == nil || value == nil {
			return false
		}
		key.(Observer).NotifyCallback(event)
		return true
	})
}
func fib(n int) chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i, j := 0, 1; i < n; i, j = i + j, i {
			out <- i
		}
	}()

	return out
}

func main() {
	n := eventSubject{
		observers: sync.Map{},
	}
	var obs1 = eventObserver{id: 1, time: time.Now()}
	var obs2 = eventObserver{id: 2, time: time.Now()}
	n.AddListener(&obs1)
	n.AddListener(&obs2)

	for x := range fib(1000000) {
		n.Nofity(Event{data : x})
	}

}
