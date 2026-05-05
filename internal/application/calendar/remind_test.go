package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/repo"
)

func TestCreateEvent_SendsToRemindChannel(t *testing.T) {
	ch := make(chan domain.Event, 10)
	svc := NewService(repo.NewEventRepoImpl(), ch)

	event := domain.Event{
		UserID:     1,
		Date:       makeDate("2099-01-01"),
		Title:      "meeting",
		RemindTime: time.Now().Add(time.Hour),
	}
	err := svc.CreateEvent(event)
	require.NoError(t, err)

	assert.Len(t, ch, 1)
	received := <-ch
	assert.Equal(t, "meeting", received.Title)
}

func TestCreateEvent_NoRemindTime_DoesNotSendToChannel(t *testing.T) {
	ch := make(chan domain.Event, 10)
	svc := NewService(repo.NewEventRepoImpl(), ch)

	err := svc.CreateEvent(domain.Event{UserID: 1, Date: makeDate("2099-01-01"), Title: "no reminder"})
	require.NoError(t, err)

	assert.Len(t, ch, 0)
}
