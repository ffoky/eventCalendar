package reminder

import (
	"eventCalendar/internal/domain"
)

type EventHeap []domain.Event

func (h EventHeap) Less(i, j int) bool {
	return h[i].RemindTime.Before(h[j].RemindTime)
}

func (h EventHeap) Len() int {
	return len(h)
}

func (h EventHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *EventHeap) Push(x any) {
	*h = append(*h, x.(domain.Event))
}

func (h *EventHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}
