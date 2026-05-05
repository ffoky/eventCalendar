package cleaner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/logger"
	"eventCalendar/internal/infrastructure/repo"
)

func TestCleanWorker_ArchivesOldEvents(t *testing.T) {
	r := repo.NewEventRepoImpl()
	require.NoError(t, r.Save(domain.Event{UserID: 1, Date: time.Now().AddDate(-1, 0, 0), Title: "old"}))
	require.NoError(t, r.Save(domain.Event{UserID: 1, Date: time.Now().AddDate(1, 0, 0), Title: "future"}))

	lgr := logger.NewLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go lgr.Run(ctx)

	w := NewCleanWorker(r, 30*time.Millisecond, lgr)
	go w.Run(ctx)

	time.Sleep(60 * time.Millisecond)

	// старое событие уже перенесено — повторный Archive вернёт 0
	assert.Equal(t, 0, r.Archive(time.Now()))

	// будущее событие на месте
	events, err := r.Get(1, time.Now().AddDate(0, -1, 0), time.Now().AddDate(2, 0, 0))
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "future", events[0].Title)
}

func TestCleanWorker_StopsOnContextCancel(t *testing.T) {
	r := repo.NewEventRepoImpl()
	lgr := logger.NewLogger()
	ctx, cancel := context.WithCancel(context.Background())

	go lgr.Run(ctx)
	w := NewCleanWorker(r, time.Hour, lgr)

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("worker did not stop after context cancel")
	}
}
