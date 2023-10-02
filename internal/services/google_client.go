package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

func NewGoogleClient(ctx context.Context, scope ...string) (*http.Client, error) {
	jsonFile, err := os.Open("google-config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	config, err := google.JWTConfigFromJSON(byteValue, scope...)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// create client with config and context
	return config.Client(ctx), nil
}