package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type GoogleCloudConfig struct {
	Type string `env:"type" json:"type"`
	ProjectId string `env:"project_id" json:"project_id"`
	PrivateKeyId string `env:"private_key_id" json:"private_key_id"`
	PrivateKey string `env:"private_key" json:"private_key"`
	ClientEmail string `env:"client_email" json:"client_email"`
	ClientId string `env:"client_id" json:"client_id"`
	AuthUri string `env:"auth_uri" json:"auth_uri"`
	TokenUri string `env:"token_uri" json:"token_uri"`
	AuthProviderX509CertUrl string `env:"auth_provider_x509_cert_url" json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl string `env:"client_x509_cert_url" json:"client_x509_cert_url"`
}

type Config struct {
	Port string `env:"PORT" json:"Port"`
	Scope string `env:"SCOPE" json:"SCOPE"`
	EventsSpreadsheetId string `env:"EVENTS_SPREADSHEET_ID" json:"EVENTS_SPREADSHEET_ID"`
	EventsReadRange string `env:"EVENTS_READ_RANGE" json:"EVENTS_READ_RANGE"`
	EntriesSpreadsheetId string `env:"ENTRIES_SPREADSHEET_ID" json:"ENTRIES_SPREADSHEET_ID"`
	EntriesReadRange string `env:"ENTRIES_READ_RANGE" json:"ENTRIES_READ_RANGE"`
	GoogleCloudConfig GoogleCloudConfig `env:"GOOGLE_CLOUD_CONFIG" json:"GOOGLE_CLOUD_CONFIG"`
}

func New() (*Config, error) {
	cfg := Config{}

	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "PORT")
	}
	log.Default().Printf("[LAN-GS-BACKEND] PORT: %v", cfg.Port)

	cfg.Scope = os.Getenv("SCOPE")
	if cfg.Scope == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "SCOPE")
	}
	log.Default().Printf("[LAN-GS-BACKEND] SCOPE: %v", cfg.Scope)

	cfg.EventsSpreadsheetId = os.Getenv("EVENTS_SPREADSHEET_ID")
	if cfg.Scope == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "EVENTS_SPREADSHEET_ID")
	}
	log.Default().Printf("[LAN-GS-BACKEND] EVENTS_SPREADSHEET_ID: %v", cfg.EventsSpreadsheetId)

	cfg.EventsReadRange = os.Getenv("EVENTS_READ_RANGE")
	if cfg.Scope == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "EVENTS_READ_RANGE")
	}
	log.Default().Printf("[LAN-GS-BACKEND] EVENTS_READ_RANGE: %v", cfg.EventsReadRange)

	cfg.EntriesSpreadsheetId = os.Getenv("ENTRIES_SPREADSHEET_ID")
	if cfg.Scope == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "ENTRIES_SPREADSHEET_ID")
	}
	log.Default().Printf("[LAN-GS-BACKEND] ENTRIES_SPREADSHEET_ID: %v", cfg.EntriesSpreadsheetId)

	cfg.EntriesReadRange = os.Getenv("ENTRIES_READ_RANGE")
	if cfg.Scope == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "ENTRIES_READ_RANGE")
	}
	log.Default().Printf("[LAN-GS-BACKEND] ENTRIES_READ_RANGE: %v", cfg.EntriesReadRange)

	googleCloudConfigString := os.Getenv("GOOGLE_CLOUD_CONFIG")
	if googleCloudConfigString == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "GOOGLE_CLOUD_CONFIG")
	}
	var googleCloudConfig GoogleCloudConfig
	if err := json.Unmarshal([]byte(googleCloudConfigString), &googleCloudConfig); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	log.Default().Printf("[LAN-GS-BACKEND] GOOGLE_CLOUD_CONFIG: %v", cfg.GoogleCloudConfig)
	cfg.GoogleCloudConfig = googleCloudConfig;

	return &cfg, nil
}
