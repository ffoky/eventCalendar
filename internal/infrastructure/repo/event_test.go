package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eventCalendar/internal/domain"
)

func makeDate(s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	return d
}

func TestArchive_MovesOldEvents(t *testing.T) {
	r := NewEventRepoImpl()
	_ = r.Save(domain.Event{UserID: 1, Date: makeDate("2020-01-01"), Title: "old"})
	_ = r.Save(domain.Event{UserID: 1, Date: makeDate("2099-01-01"), Title: "future"})

	count := r.Archive(time.Now())

	assert.Equal(t, 1, count)

	events, err := r.Get(1, makeDate("2020-01-01"), makeDate("2020-01-01"))
	require.NoError(t, err)
	assert.Empty(t, events)

	events, err = r.Get(1, makeDate("2099-01-01"), makeDate("2099-01-01"))
	require.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestArchive_Empty(t *testing.T) {
	r := NewEventRepoImpl()
	count := r.Archive(time.Now())
	assert.Equal(t, 0, count)
}
