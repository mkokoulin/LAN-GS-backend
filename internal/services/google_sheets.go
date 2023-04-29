package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	spreadsheetId string
	readRange string
	srv *sheets.Service
}

type Event struct {
	Name string `mapstructure:"name"`	
	Description	string `mapstructure:"description"`
	DateFrom string `mapstructure:"date_from"`
	DateTo string `mapstructure:"date_to"`
	GoogleForm string `mapstructure:"google_form"`
	Payment	string `mapstructure:"payment"`
}

func NewGoogleSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*SheetsService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &SheetsService {
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (SS *SheetsService) GetEvents(ctx context.Context) ([]Event, error) {
	res, err := SS.srv.Spreadsheets.Values.Get(SS.spreadsheetId, SS.readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("%v", err)
	}

	events := []Event{}

	colMap := map[int]string {
		0: "name",
		1: "description",
		2: "date_from",
		3: "date_to",
		4: "google_form",
		5: "payment",
	}

	for _, val := range res.Values {
		e := map[string]interface{} {}

		for i, v := range val {
			col, ok := colMap[i]
			if ok {
				e[col] = v.(string)
			}
		}

		var event Event

		mapstructure.Decode(e, &event)
		
		events = append(events, event)
	}

	return events, nil
}