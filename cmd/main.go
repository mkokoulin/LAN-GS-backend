package main

import (
	"context"
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

	gc, err := services.NewGoogleClient(ctx, cfg.Google.SecretPath, cfg.Google.Scope)
	if err != nil {
		log.Fatal(err)
		return
	}

	ss, err := services.NewGoogleSheets(ctx, gc, cfg.GoogleSheets.SpreadsheetId, cfg.GoogleSheets.ReadRange)
	if err != nil {
		log.Fatal(err)
		return
	}

	h := handlers.New(ss)

	r := router.New(h)

	s := server.New(cfg.ServerAddress, r)

	g.Go(func() error {
		s.Start()

		log.Printf("httpServer starting at: %v", cfg.ServerAddress)

		return nil
	})

	select {
	case <-interrupt:
		log.Println("Stop server")
		break
	case <-ctx.Done():
		break
	}

	err = g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
	}
}