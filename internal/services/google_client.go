package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mkokoulin/LAN-GS-backend/internal/config"
	"golang.org/x/oauth2/google"
)

func NewGoogleClient(ctx context.Context, gcc config.GoogleCloudConfig, scope ...string) (*http.Client, error) {
	byteValue, err := json.Marshal(gcc)
	if err != nil {
		fmt.Println(err)
	}

	config, err := google.JWTConfigFromJSON(byteValue, scope...)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return config.Client(ctx), nil
}