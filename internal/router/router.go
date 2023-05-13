package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mkokoulin/LAN-GS-backend/internal/handlers"
)

// New router constructor
func New(h *handlers.Handlers) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/api/events", h.GetEvents)
		r.Post("/api/events/update", h.UpdateEvent)
		r.Post("/api/entries", h.CreateEntrie)
		r.Post("/api/entries/cancel", h.CancelEntrie)
		r.Put("/api/entries/update", h.UpdateEntrie)
	})

	return router
}