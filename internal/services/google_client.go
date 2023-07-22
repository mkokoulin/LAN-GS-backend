package services

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"
)

func NewGoogleClient(ctx context.Context, secret string, scope ...string) (*http.Client, error) {
	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON([]byte(secret), scope...)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// create client with config and context
	return config.Client(ctx), nil
}