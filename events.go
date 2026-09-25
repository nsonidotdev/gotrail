package gotrail

import "time"

type Event struct {
	Name       string
	StartTime  time.Time
	Attributes map[string]any
}

func (s *Span) AddEvent(name string, attributes map[string]any) {
	event := &Event{
		name,
		time.Now(),
		attributes,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Events = append(s.Events, event)
}
