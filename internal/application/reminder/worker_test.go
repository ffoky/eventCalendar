package reminder

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/logger"
)

func TestRemindWorker_FiresPastReminder(t *testing.T) {
	ch := make(chan domain.Event, 1)
	lgr := logger.NewLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go lgr.Run(ctx)

	w := NewRemindWorker(ch, lgr)

	notified := make(chan domain.Event, 1)
	w.notify = func(e domain.Event) { notified <- e }

	go w.Run(ctx)

	ch <- domain.Event{
		UserID:     1,
		Title:      "past reminder",
		RemindTime: time.Now().Add(-time.Second),
	}

	select {
	case e := <-notified:
		assert.Equal(t, "past reminder", e.Title)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("reminder did not fire in time")
	}
}

func TestRemindWorker_FiresInOrder(t *testing.T) {
	ch := make(chan domain.Event, 3)
	lgr := logger.NewLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go lgr.Run(ctx)

	w := NewRemindWorker(ch, lgr)

	notified := make(chan string, 3)
	w.notify = func(e domain.Event) { notified <- e.Title }

	go w.Run(ctx)

	now := time.Now()
	// отправляем в "неправильном" порядке — heap должен исправить
	ch <- domain.Event{Title: "third", RemindTime: now.Add(30 * time.Millisecond)}
	ch <- domain.Event{Title: "first", RemindTime: now.Add(10 * time.Millisecond)}
	ch <- domain.Event{Title: "second", RemindTime: now.Add(20 * time.Millisecond)}

	var order []string
	for range 3 {
		select {
		case title := <-notified:
			order = append(order, title)
		case <-time.After(500 * time.Millisecond):
			t.Fatal("reminder did not fire in time")
		}
	}
	assert.Equal(t, []string{"first", "second", "third"}, order)
}
