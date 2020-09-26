package internalhttp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

const timeout = time.Second

var errWrongPeriod = errors.New("wrong period")

func getPathAndID(urlPath string) (string, int) {
	var path string
	var id int

	lastInd := strings.LastIndex(urlPath, "/")
	itemID := urlPath[lastInd+1:]
	id, err := strconv.Atoi(itemID)
	if err == nil {
		path = urlPath[:lastInd]
	} else {
		path = urlPath
	}

	return path, id
}

// EventHandler ...
type EventHandler struct {
	app app.Application
}

// NewEventHandler ...
func NewEventHandler(a app.Application) *EventHandler {
	return &EventHandler{app: a}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	path, id := getPathAndID(r.URL.Path)
	if path != "/event" {
		http.NotFound(w, r)

		return
	}
	if id != 0 {
		switch r.Method {
		case "GET":
			h.GetEvent(ctx, w, r)
		case "PATCH":
			h.UpdateEvent(ctx, w, r)
		case "DELETE":
			h.DeleteEvent(ctx, w, r)
		default:
		}
	} else {
		switch r.Method {
		case "GET":
			h.GetEventList(ctx, w, r)
		case "POST":
			h.CreateEvent(ctx, w, r)
		default:
		}
	}
}

// GetEvent returns event by id.
func (h *EventHandler) GetEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, id := getPathAndID(r.URL.Path)
	e, err := h.app.GetEvent(ctx, int64(id))
	switch err {
	case nil:
	case app.ErrEventIDZero:
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	case sql.ErrNoRows:
		http.Error(w, "event not found", http.StatusNotFound)

		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	js, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) //nolint
}

// GetEventList returns list of events.
func (h *EventHandler) GetEventList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var rawDate string
	var rawPeriod string
	for qName, qVal := range r.URL.Query() {
		if qName == "date" {
			rawDate = qVal[0]
		}
		if qName == "period" {
			rawPeriod = qVal[0]
		}
	}

	if len(rawDate) == 0 || len(rawPeriod) == 0 {
		http.Error(w, "required query params missed", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", rawDate)
	if err != nil {
		http.Error(w, "date wrong format. required [2006-01-02]", http.StatusBadRequest)
	}

	e, err := h.getList(ctx, date, rawPeriod)
	switch err {
	case nil:
	case errWrongPeriod:
		http.Error(w, "wrong period. allowed d/w/m", http.StatusBadRequest)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) //nolint
}

func (h *EventHandler) getList(ctx context.Context, date time.Time, period string) ([]*storage.Event, error) {
	switch period {
	case "d":
		return h.app.GetEventForDay(ctx, date)
	case "w":
		return h.app.GetEventForWeek(ctx, date)
	case "m":
		return h.app.GetEventForMonth(ctx, date)
	default:
		return nil, errWrongPeriod
	}
}

// CreateEvent creates new event.
func (h *EventHandler) CreateEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)

		return
	}
	e := &storage.Event{}
	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	e, err = h.app.CreateEvent(ctx,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.OwnerID,
		e.Descr,
		e.NotifyBefore,
	)
	switch err {
	case nil:
	case app.ErrEventFieldWrong:
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	js, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(js) //nolint
}

// UpdateEvent updates event.
func (h *EventHandler) UpdateEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, id := getPathAndID(r.URL.Path)
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)

		return
	}
	e := &storage.Event{}
	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	e, err = h.app.UpdateEvent(ctx,
		int64(id),
		e.Title,
		e.StartDate,
		e.EndDate,
		e.OwnerID,
		e.Descr,
		e.NotifyBefore,
	)
	switch err {
	case nil:
	case app.ErrEventIDZero:
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	js, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) //nolint
}

// DeleteEvent deletes event.
func (h *EventHandler) DeleteEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, id := getPathAndID(r.URL.Path)
	err := h.app.DeleteEvent(ctx, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
