package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/repo"
)

func newTestService() *Service {
	return NewService(repo.NewEventRepoImpl(), make(chan domain.Event, 100))
}

func makeDate(s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	return d
}

func TestCreateAndGetForDay(t *testing.T) {
	s := newTestService()

	event := domain.Event{UserID: 1, Date: makeDate("2026-03-21"), Title: "meeting"}
	err := s.CreateEvent(event)
	require.NoError(t, err)

	events, err := s.GetForDay(1, makeDate("2026-03-21"))
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "meeting", events[0].Title)
}

func TestGetForDayEmpty(t *testing.T) {
	s := newTestService()

	events, err := s.GetForDay(1, makeDate("2026-03-21"))
	require.NoError(t, err)
	assert.Nil(t, events)
}

func TestGetForWeek(t *testing.T) {
	s := newTestService()

	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-16"), Title: "monday"})
	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-18"), Title: "wednesday"})
	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-22"), Title: "sunday"})
	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-23"), Title: "next week"})

	events, err := s.GetForWeek(1, makeDate("2026-03-19"))
	require.NoError(t, err)
	assert.Len(t, events, 3)
}

func TestGetForMonth(t *testing.T) {
	s := newTestService()

	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-01"), Title: "first"})
	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-31"), Title: "last"})
	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-04-01"), Title: "next month"})

	events, err := s.GetForMonth(1, makeDate("2026-03-15"))
	require.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestUpdateEvent(t *testing.T) {
	s := newTestService()

	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-21"), Title: "old"})

	events, _ := s.GetForDay(1, makeDate("2026-03-21"))
	updated := events[0]
	updated.Title = "new"
	err := s.UpdateEvent(updated)
	require.NoError(t, err)

	events, _ = s.GetForDay(1, makeDate("2026-03-21"))
	assert.Equal(t, "new", events[0].Title)
}

func TestUpdateEventNotFound(t *testing.T) {
	s := newTestService()

	err := s.UpdateEvent(domain.Event{ID: 999, UserID: 1, Date: makeDate("2026-03-21"), Title: "x"})
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestDeleteEvent(t *testing.T) {
	s := newTestService()

	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-21"), Title: "to delete"})

	events, _ := s.GetForDay(1, makeDate("2026-03-21"))
	err := s.DeleteEvent(events[0].ID, 1)
	require.NoError(t, err)

	events, _ = s.GetForDay(1, makeDate("2026-03-21"))
	assert.Nil(t, events)
}

func TestDeleteEventNotFound(t *testing.T) {
	s := newTestService()

	err := s.DeleteEvent(999, 1)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestDifferentUsers(t *testing.T) {
	s := newTestService()

	_ = s.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2026-03-21"), Title: "user1"})
	_ = s.CreateEvent(domain.Event{UserID: 2, Date: makeDate("2026-03-21"), Title: "user2"})

	events, _ := s.GetForDay(1, makeDate("2026-03-21"))
	assert.Len(t, events, 1)
	assert.Equal(t, "user1", events[0].Title)
}
