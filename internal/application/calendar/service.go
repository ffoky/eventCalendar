package calendar

import (
	"time"

	"eventCalendar/internal/domain"
)

type Service struct {
	repo      domain.EventRepo
	reminders chan<- domain.Event
}

func NewService(repo domain.EventRepo, reminders chan<- domain.Event) *Service {
	return &Service{repo: repo, reminders: reminders}
}

func (s *Service) CreateEvent(event domain.Event) error {
	if err := s.repo.Save(event); err != nil {
		return err
	}
	if !event.RemindTime.IsZero() {
		select {
		case s.reminders <- event:
		default:
		}
	}
	return nil
}

func (s *Service) UpdateEvent(event domain.Event) error {
	return s.repo.Update(event)
}

func (s *Service) DeleteEvent(ID, userID int) error {
	return s.repo.Delete(ID, userID)
}

func (s *Service) GetForDay(userID int, date time.Time) ([]domain.Event, error) {
	return s.repo.Get(userID, date, date)
}

func (s *Service) GetForWeek(userID int, date time.Time) ([]domain.Event, error) {
	dateWeekDay := date.Weekday()
	if dateWeekDay > 0 {
		dateWeekDay -= 1
	} else {
		dateWeekDay = 6
	}
	from := date.AddDate(0, 0, -int(dateWeekDay))
	to := from.AddDate(0, 0, 6)
	return s.repo.Get(userID, from, to)
}

func (s *Service) GetForMonth(userID int, date time.Time) ([]domain.Event, error) {
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	to := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
	return s.repo.Get(userID, from, to)
}
