package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eventCalendar/internal/application/calendar"
	"eventCalendar/internal/domain"
	"eventCalendar/internal/infrastructure/repo"
)

func newTestHandler() *EventHandler {
	ch := make(chan domain.Event, 10)
	svc := calendar.NewService(repo.NewEventRepoImpl(), ch)
	return NewEventHandler(svc)
}

func TestCreateEvent_OK(t *testing.T) {
	h := newTestHandler()
	body := `{"user_id":1,"date":"2099-01-01","title":"meeting"}`
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp successResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "event created", resp.Result)
}

func TestCreateEvent_BadDate(t *testing.T) {
	h := newTestHandler()
	body := `{"user_id":1,"date":"not-a-date","title":"meeting"}`
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteEvent_NotFound(t *testing.T) {
	h := newTestHandler()
	body := `{"id":999,"user_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/delete_event", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestEventsForDay_BadParams(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/events_for_day?user_id=abc&date=2099-01-01", nil)
	w := httptest.NewRecorder()

	h.Mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
