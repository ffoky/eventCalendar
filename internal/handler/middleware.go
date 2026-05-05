package handler

import (
	"fmt"
	"net/http"
	"time"

	"eventCalendar/internal/infrastructure/logger"
)

func LogMiddleware(logger *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Log(fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, time.Since(start)))
	})
}
