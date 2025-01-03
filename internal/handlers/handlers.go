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
	GetEvent(ctx context.Context, eventId string) (services.Event, error)
	UpdateEvent(ctx context.Context, event services.Event) error
}

type GSEntriesInterface interface {
	CreateEntrie(ctx context.Context, entrie services.Entrie) (services.EntrieResponse, error)
	UpdateEntrie(ctx context.Context, entrie services.Entrie) error
	CancelEntrie(ctx context.Context, entrie services.CancelEntrieResponse) error
	GetUniqueEntries(ctx context.Context) (map[string]int, error)
}

type B2bRequestsInterface interface {
	CreateB2bRequest(ctx context.Context, b2bRequest services.B2bRequest) error
}

type BookingInterface interface {
	CreateBooking(ctx context.Context, booking services.Booking) error
}

type Handlers struct {
	events GSEventsInterface
	entries GSEntriesInterface
	b2bRequests B2bRequestsInterface
	booking BookingInterface
}

func New(events GSEventsInterface, entries GSEntriesInterface, b2bRequests B2bRequestsInterface, booking BookingInterface) *Handlers {
	return &Handlers{
		events: events,
		entries: entries,
		b2bRequests: b2bRequests,
		booking: booking,
	}
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
    (*w).Header().Set("Access-Control-Allow-Origin", "https://lettersandnumbers.am")
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

func (h *Handlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)

	// s := r.Header.Get("X-Secret")

	// if s != "0a41c238c148d57ab77a850ce491bc00" {
	// 	w.WriteHeader(http.StatusNotFound);
	// 	return;
	// }

	eventId := r.URL.Query().Get("eventId")

	if eventId == "" {
		http.Error(w, "eventId cannot be an empty", http.StatusBadRequest)
		return
	}

	event, err := h.events.GetEvent(r.Context(), eventId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if event.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(event)

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

	// s := r.Header.Get("X-Secret")

	// if s != "0a41c238c148d57ab77a850ce491bc00" {
	// 	w.WriteHeader(http.StatusNotFound);
	// 	return;
	// }

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

func (h *Handlers) CreateBooking(w http.ResponseWriter, r *http.Request) {
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

	booking := services.Booking {}

	err = json.Unmarshal(body, &booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.booking.CreateBooking(r.Context(), booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UpdateEntrie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	setupCORS(&w, r)

	// s := r.Header.Get("X-Secret")

	// if s != "0a41c238c148d57ab77a850ce491bc00" {
	// 	w.WriteHeader(http.StatusNotFound);
	// 	return;
	// }

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

	// s := r.Header.Get("X-Secret")

	// if s != "0a41c238c148d57ab77a850ce491bc00" {
	// 	w.WriteHeader(http.StatusNotFound);
	// 	return;
	// }

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

	// s := r.Header.Get("X-Secret")

	// if s != "0a41c238c148d57ab77a850ce491bc00" {
	// 	w.WriteHeader(http.StatusNotFound);
	// 	return;
	// }

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

func (h *Handlers) GetUniqueEntries(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)

	// s := r.Header.Get("X-Secret")

	// if s != "0a41c238c148d57ab77a850ce491bc00" {
	// 	w.WriteHeader(http.StatusNotFound);
	// 	return;
	// }

	entries, err := h.entries.GetUniqueEntries(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if len(entries) == 0 {
	// 	http.Error(w, errors.New("no content").Error(), http.StatusNoContent)
	// 	return
	// }

	body, err := json.Marshal(entries)

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(body)
		if err == nil {
			return
		}
	}
}

func (h *Handlers) CreateB2bRequest(w http.ResponseWriter, r *http.Request) {
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

	b2bRequest := services.B2bRequest {}

	err = json.Unmarshal(body, &b2bRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.b2bRequests.CreateB2bRequest(r.Context(), b2bRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)
	}
}
