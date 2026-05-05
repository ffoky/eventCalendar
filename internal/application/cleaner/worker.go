package cleaner

import (
	"context"
	"time"

	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/logger"
)

type CleanWorker struct {
	repo     domain.EventRepo
	interval time.Duration
	logger   *logger.Logger
}

func NewCleanWorker(repo domain.EventRepo, interval time.Duration, logger *logger.Logger) *CleanWorker {
	return &CleanWorker{
		repo:     repo,
		interval: interval,
		logger:   logger,
	}
}

func (w *CleanWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			archived := w.repo.Archive(time.Now())
			w.logger.Logf("archived: %d events", archived)
		case <-ctx.Done():
			return
		}
	}
}
