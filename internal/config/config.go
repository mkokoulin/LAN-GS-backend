package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type GoogleConfig struct {
	Scope string `env:"SCOPE" json:"SCOPE"`
	GoogleSecret string `env:"GOOGLE_SECRET" json:"GOOGLE_SECRET"`
}

type EventTable struct {
	SpreadsheetId string `env:"EVENT_SPREADSHEET_ID" json:"EVENT_SPREADSHEET_ID"`
	ReadRange string `env:"EVENT_READ_RANGE" json:"EVENT_READ_RANGE"`
}

type EntriesTable struct {
	SpreadsheetId string `env:"ENTRIES_SPREADSHEET_ID" json:"ENTRIES_SPREADSHEET_ID"`
	ReadRange string `env:"ENTRIES_READ_RANGE" json:"ENTRIES_READ_RANGE"`
}

type GoogleSheetsConfig struct {
	Event EventTable
	Entries EntriesTable
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

	eventTable := EventTable{}
	entriesTable := EntriesTable{}

	err = env.Parse(&eventTable)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	err = env.Parse(&entriesTable)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	gsc.Event = eventTable
	gsc.Entries = entriesTable

	cfg.Google = gc
	cfg.GoogleSheets = gsc

	return &cfg, nil
}