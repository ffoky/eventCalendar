package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"eventCalendar/internal/application/calendar"
	"eventCalendar/internal/domain"
)

var (
	ErrInvalidUserID = errors.New("invalid user_id")
	ErrInvalidDate   = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrInvalidBody   = errors.New("invalid request body")
)

type createEventRequest struct {
	UserID     int    `json:"user_id"`
	Date       string `json:"date"`
	Title      string `json:"title"`
	RemindTime string `json:"remind_time"`
}

type updateEventRequest struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Date       string `json:"date"`
	Title      string `json:"title"`
	RemindTime string `json:"remind_time"`
}

type deleteEventRequest struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

type successResponse struct {
	Result string `json:"result"`
}

type eventResponse struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Date       string `json:"date"`
	Title      string `json:"title"`
	RemindTime string `json:"remind_time,omitempty"`
}

type eventsResponse struct {
	Result []eventResponse `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type EventHandler struct {
	service *calendar.Service
	Mux     *http.ServeMux
}

func NewEventHandler(service *calendar.Service) *EventHandler {
	h := &EventHandler{
		service: service,
		Mux:     http.NewServeMux(),
	}
	h.registerRoutes()
	return h
}

func (h *EventHandler) registerRoutes() {
	h.Mux.HandleFunc("POST /create_event", h.createEvent)
	h.Mux.HandleFunc("POST /update_event", h.updateEvent)
	h.Mux.HandleFunc("POST /delete_event", h.deleteEvent)
	h.Mux.HandleFunc("GET /events_for_day", h.eventsForDay)
	h.Mux.HandleFunc("GET /events_for_week", h.eventsForWeek)
	h.Mux.HandleFunc("GET /events_for_month", h.eventsForMonth)
}

func (h *EventHandler) createEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidBody.Error())
		return
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidDate.Error())
		return
	}
	var remind time.Time
	if req.RemindTime != "" {
		remind, err = time.Parse("2006-01-02T15:04:05", req.RemindTime)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid remind_time format, expected YYYY-MM-DDTHH:MM:SS")
			return
		}
	}
	event := domain.Event{
		UserID:     req.UserID,
		Date:       date,
		Title:      req.Title,
		RemindTime: remind,
	}
	if err = h.service.CreateEvent(event); err != nil {
		writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Result: "event created"})
}

func (h *EventHandler) updateEvent(w http.ResponseWriter, r *http.Request) {
	var req updateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidBody.Error())
		return
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidDate.Error())
		return
	}
	var remind time.Time
	if req.RemindTime != "" {
		remind, err = time.Parse("2006-01-02T15:04:05", req.RemindTime)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid remind_time format, expected YYYY-MM-DDTHH:MM:SS")
			return
		}
	}
	event := domain.Event{
		ID:         req.ID,
		UserID:     req.UserID,
		Date:       date,
		Title:      req.Title,
		RemindTime: remind,
	}
	if err = h.service.UpdateEvent(event); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			writeError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Result: "event updated"})
}

func (h *EventHandler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	var req deleteEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidBody.Error())
		return
	}
	if err := h.service.DeleteEvent(req.ID, req.UserID); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			writeError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Result: "event deleted"})
}

func (h *EventHandler) eventsForDay(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseGetParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := h.service.GetForDay(userID, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, eventsResponse{Result: toEventResponses(events)})
}

func (h *EventHandler) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseGetParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := h.service.GetForWeek(userID, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, eventsResponse{Result: toEventResponses(events)})
}

func (h *EventHandler) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseGetParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := h.service.GetForMonth(userID, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, eventsResponse{Result: toEventResponses(events)})
}

func parseGetParams(r *http.Request) (int, time.Time, error) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return 0, time.Time{}, ErrInvalidUserID
	}
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		return 0, time.Time{}, ErrInvalidDate
	}
	return userID, date, nil
}

func toEventResponses(events []domain.Event) []eventResponse {
	result := make([]eventResponse, len(events))
	for i, e := range events {
		resp := eventResponse{
			ID:     e.ID,
			UserID: e.UserID,
			Date:   e.Date.Format("2006-01-02"),
			Title:  e.Title,
		}
		if !e.RemindTime.IsZero() {
			resp.RemindTime = e.RemindTime.Format("2006-01-02T15:04:05")
		}
		result[i] = resp
	}
	return result
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
