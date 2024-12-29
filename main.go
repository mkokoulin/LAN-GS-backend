package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mkokoulin/LAN-GS-backend/internal/config"
	"github.com/mkokoulin/LAN-GS-backend/internal/handlers"
	"github.com/mkokoulin/LAN-GS-backend/internal/router"
	"github.com/mkokoulin/LAN-GS-backend/internal/server"
	"github.com/mkokoulin/LAN-GS-backend/internal/services"
	"golang.org/x/sync/errgroup"
)

func main() {
	// create api context
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)
	interrupt := make(chan os.Signal, 1)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	gc, err := services.NewGoogleClient(ctx, cfg.GoogleCloudConfig, cfg.Scope)
	if err != nil {
		log.Fatal(err)
		return
	}

	eventsSheets, err := services.NewEventsSheets(ctx, gc, cfg.EventsSpreadsheetId, cfg.EventsReadRange)
	if err != nil {
		log.Fatal(err)
		return
	}

	entriesSheets, err := services.NewEntriesSheets(ctx, gc, cfg.EntriesSpreadsheetId, cfg.EntriesReadRange)
	if err != nil {
		log.Fatal(err)
		return
	}

	b2bRequestSheets, err := services.NewB2bRequestSheets(ctx, gc, cfg.B2bRequestsSpreadsheetId, cfg.EntriesReadRange)
	if err != nil {
		log.Fatal(err)
		return
	}

	bookingSheets, err := services.NewBookingSheets(ctx, gc, cfg.BookingSpreadsheetId)
	if err != nil {
		log.Fatal(err)
		return
	}

	h := handlers.New(eventsSheets, entriesSheets, b2bRequestSheets, bookingSheets)

	r := router.New(h)

	s := server.New(r, cfg.Port)

	g.Go(func() error {
		_, err = s.Start()
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		return nil
	})

	select {
	case <-interrupt:
		// stop(ctx)
		log.Println("Stop server")
		break
	case <-ctx.Done():
		break
	}

	err = g.Wait()
	if err != nil {
		// stop(ctx)
		log.Printf("server returning an error: %v", err)
		return
	}
}