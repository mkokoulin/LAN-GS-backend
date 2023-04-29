package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mkokoulin/LAN-GS-backend/internal/services"
)

type GSServiceInterface interface {
	GetEvents(ctx context.Context) ([]services.Event, error)
}

type Handlers struct {
	service GSServiceInterface
}

func New(service GSServiceInterface) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) GetEvents(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		http.Error(w, errors.New("no content").Error(), http.StatusNoContent)
		return
	}

	body, err := json.Marshal(urls)

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(body)
		if err == nil {
			return
		}
	}
}