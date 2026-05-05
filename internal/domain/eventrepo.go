package domain

import "time"

type EventRepo interface {
	Save(event Event) error
	Update(event Event) error
	Delete(ID, userID int) error
	Get(userID int, from, to time.Time) ([]Event, error)
	Archive(before time.Time) int
}
