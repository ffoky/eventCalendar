package repo

import (
	"fmt"
	"sync"
	"time"

	"eventCalendar/internal/domain"
)

type EventRepoImpl struct {
	vault   map[int]domain.Event
	archive map[int]domain.Event
	mu      sync.RWMutex
	nextID  int
}

func NewEventRepoImpl() *EventRepoImpl {
	return &EventRepoImpl{
		vault:   make(map[int]domain.Event),
		archive: make(map[int]domain.Event),
	}
}

func (er *EventRepoImpl) Save(event domain.Event) error {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.nextID++
	event.ID = er.nextID
	er.vault[event.ID] = event
	return nil
}

func (er *EventRepoImpl) Update(event domain.Event) error {
	er.mu.Lock()
	defer er.mu.Unlock()
	if _, ok := er.vault[event.ID]; !ok {
		return fmt.Errorf("update event in repo: %w", domain.ErrEventNotFound)
	}
	er.vault[event.ID] = event
	return nil
}

func (er *EventRepoImpl) Get(userID int, from, to time.Time) ([]domain.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()
	res := make([]domain.Event, 0)
	for _, event := range er.vault {
		if event.UserID == userID && !event.Date.Before(from) && !event.Date.After(to) {
			res = append(res, event)
		}
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (er *EventRepoImpl) Delete(ID, userID int) error {
	er.mu.Lock()
	defer er.mu.Unlock()
	event, ok := er.vault[ID]
	if !ok {
		return fmt.Errorf("delete event in repo: %w", domain.ErrEventNotFound)
	}
	if event.UserID != userID {
		return fmt.Errorf("delete event in repo: %w", domain.ErrEventNotFound)
	}
	delete(er.vault, ID)
	return nil
}

func (er *EventRepoImpl) Archive(before time.Time) int {
	er.mu.Lock()
	defer er.mu.Unlock()
	count := 0
	for _, event := range er.vault {
		if event.Date.Before(before) {
			count++
			delete(er.vault, event.ID)
			er.archive[event.ID] = event
		}
	}
	return count
}
