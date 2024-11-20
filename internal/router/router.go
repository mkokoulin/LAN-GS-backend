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
		r.Get("/api/event", h.GetEvent)
		r.Post("/api/events/update", h.UpdateEvent)
		r.Get("/api/entries/unique", h.GetUniqueEntries)
		r.Post("/api/entries", h.CreateEntrie)
		r.Post("/api/entries/cancel", h.CancelEntrie)
		r.Put("/api/entries/update", h.UpdateEntrie)
		r.Post("/api/b2b/request", h.CreateB2bRequest)
	})

	return router
}