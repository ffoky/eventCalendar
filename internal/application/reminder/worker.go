package reminder

import (
	"container/heap"
	"context"
	"time"

	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/logger"
)

type RemindWorker struct {
	events <-chan domain.Event
	logger *logger.Logger
	notify func(domain.Event)
}

func NewRemindWorker(ch <-chan domain.Event, lgr *logger.Logger) *RemindWorker {
	w := &RemindWorker{events: ch, logger: lgr}
	w.notify = w.remind
	return w
}

func (w *RemindWorker) Run(ctx context.Context) {
	h := EventHeap{}
	heap.Init(&h)
	timer := time.NewTimer(0)
	<-timer.C

	for {
		select {
		case event := <-w.events:
			heap.Push(&h, event)
			resetTimer(timer, h[0].RemindTime)
		case <-timer.C:
			if h.Len() == 0 {
				continue
			}
			e := heap.Pop(&h).(domain.Event)
			w.notify(e)
			if h.Len() > 0 {
				resetTimer(timer, h[0].RemindTime)
			}
		case <-ctx.Done():
			return
		}
	}
}

func resetTimer(timer *time.Timer, at time.Time) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(time.Until(at))
}

func (w *RemindWorker) remind(event domain.Event) {
	w.logger.Logf("reminder: user=%d title=%q date=%s", event.UserID, event.Title, event.Date.Format("2006-01-02"))
}
