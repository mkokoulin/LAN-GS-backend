package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mkokoulin/LAN-GS-backend/internal/services"
)

type GSEventsInterface interface {
	GetEvents(ctx context.Context) ([]services.Event, error)
	UpdateEvent(ctx context.Context, event services.Event) error
}

type GSEntriesInterface interface {
	CreateEntrie(ctx context.Context, entrie services.Entrie) (services.EntrieResponse, error)
	UpdateEntrie(ctx context.Context, entrie services.Entrie) error
	CancelEntrie(ctx context.Context, entrie services.CancelEntrieResponse) error
}

type Handlers struct {
	events GSEventsInterface
	entries GSEntriesInterface
}

func New(events GSEventsInterface, entries GSEntriesInterface) *Handlers {
	return &Handlers{
		events: events,
		entries: entries,
	}
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (h *Handlers) GetEvents(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)

	events, err := h.events.GetEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(events) == 0 {
		http.Error(w, errors.New("no content").Error(), http.StatusNoContent)
		return
	}

	body, err := json.Marshal(events)

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(body)
		if err == nil {
			return
		}
	}
}

func (h *Handlers) CreateEntrie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	setupCORS(&w, r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	entrie := services.Entrie {}

	err = json.Unmarshal(body, &entrie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entrieRes, err := h.entries.CreateEntrie(r.Context(), entrie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err = json.Marshal(entrieRes)

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(body)
		if err == nil {
			return
		}
	}
}

func (h *Handlers) UpdateEntrie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// setupCORS(&w, r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	entrie := services.Entrie {}

	err = json.Unmarshal(body, &entrie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.entries.UpdateEntrie(r.Context(), entrie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {	
	defer r.Body.Close()

	setupCORS(&w, r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	event := services.Event {}

	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.events.UpdateEvent(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) CancelEntrie(w http.ResponseWriter, r *http.Request) {	
	defer r.Body.Close()

	setupCORS(&w, r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "the body cannot be an empty", http.StatusBadRequest)
		return
	}

	event := services.CancelEntrieResponse {}

	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.entries.CancelEntrie(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusOK)
	}
}