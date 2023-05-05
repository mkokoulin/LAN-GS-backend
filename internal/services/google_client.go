package services

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

func NewGoogleClient(ctx context.Context, secretPath string, scope ...string) (*http.Client, error) {
	data, err := os.ReadFile(secretPath)
    if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
    }
	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(data, scope...)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// create client with config and context
	return config.Client(ctx), nil
}