package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eventCalendar/internal/application/calendar"
	"eventCalendar/internal/application/cleaner"
	"eventCalendar/internal/application/reminder"
	"eventCalendar/internal/domain"
	"eventCalendar/internal/handler"
	"eventCalendar/internal/infrastructure/logger"
	"eventCalendar/internal/infrastructure/repo"
)

const (
	defaultPort          = 8080
	defaultCleanInterval = 5 * time.Minute
	shutdownTimeout      = 5 * time.Second
	remindChSize         = 100
)

func main() {
	port := flag.Int("port", defaultPort, "server port")
	cleanInterval := flag.Duration("clean-interval", defaultCleanInterval, "interval for archiving old events")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	lgr := logger.NewLogger()
	go lgr.Run(ctx)

	eventRepo := repo.NewEventRepoImpl()
	remindCh := make(chan domain.Event, remindChSize)

	svc := calendar.NewService(eventRepo, remindCh)
	h := handler.NewEventHandler(svc)

	remindWorker := reminder.NewRemindWorker(remindCh, lgr)
	cleanWorker := cleaner.NewCleanWorker(eventRepo, *cleanInterval, lgr)

	go remindWorker.Run(ctx)
	go cleanWorker.Run(ctx)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler.LogMiddleware(lgr, h.Mux),
	}

	lgr.Logf("starting server on :%d", *port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		lgr.Logf("shutdown error: %s", err)
	}
}
