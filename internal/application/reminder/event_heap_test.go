package reminder

import (
	"container/heap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"eventCalendar/internal/domain"
)

func TestEventHeap_Order(t *testing.T) {
	h := &EventHeap{}
	heap.Init(h)

	heap.Push(h, domain.Event{Title: "third", RemindTime: time.Now().Add(3 * time.Hour)})
	heap.Push(h, domain.Event{Title: "first", RemindTime: time.Now().Add(1 * time.Hour)})
	heap.Push(h, domain.Event{Title: "second", RemindTime: time.Now().Add(2 * time.Hour)})

	assert.Equal(t, "first", heap.Pop(h).(domain.Event).Title)
	assert.Equal(t, "second", heap.Pop(h).(domain.Event).Title)
	assert.Equal(t, "third", heap.Pop(h).(domain.Event).Title)
}
