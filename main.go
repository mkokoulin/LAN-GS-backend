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
	// helpers.Generate()
	// create api context
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)
	interrupt := make(chan os.Signal, 1)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	gc, err := services.NewGoogleClient(ctx, cfg.Google.GoogleSecret, cfg.Google.Scope)
	if err != nil {
		log.Fatal(err)
		return
	}

	eventsSheets, err := services.NewEventsSheets(ctx, gc, cfg.GoogleSheets.Event.SpreadsheetId, cfg.GoogleSheets.Event.ReadRange)
	if err != nil {
		log.Fatal(err)
		return
	}

	entriesSheets, err := services.NewEntriesSheets(ctx, gc, cfg.GoogleSheets.Entries.SpreadsheetId, cfg.GoogleSheets.Entries.ReadRange)
	if err != nil {
		log.Fatal(err)
		return
	}

	h := handlers.New(eventsSheets, entriesSheets)

	r := router.New(h)

	s := server.New(r)

	// var stop func(ctx context.Context) error

	g.Go(func() error {
		// err = s.StartTLS("taplink-cert.pem", "taplink-key.pem")
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