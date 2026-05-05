package domain

import (
	"errors"
	"time"
)

var (
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventNotFound      = errors.New("event not found")
)

type Event struct {
	ID         int
	UserID     int
	Date       time.Time
	Title      string
	RemindTime time.Time
}
