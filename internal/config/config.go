package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type GoogleConfig struct {
	SecretPath string `env:"SECRET_PATH" json:"SECRET_PATH"`
	Scope string `env:"SCOPE" json:"SCOPE"`
}

type GoogleSheetsConfig struct {
	SpreadsheetId string `env:"SPREADSHEET_ID" json:"SPREADSHEET_ID"`
	ReadRange string `env:"READ_RANGE" json:"READ_RANGE"`
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" json:"SERVER_ADDRESS"`
	Google GoogleConfig
	GoogleSheets GoogleSheetsConfig
}

func New() (*Config, error) {
	cfg := Config{}
	
	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	gc := GoogleConfig{}

	err = env.Parse(&gc)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	gsc := GoogleSheetsConfig{}

	err = env.Parse(&gsc)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	cfg.Google = gc
	cfg.GoogleSheets = gsc

	return &cfg, nil
}