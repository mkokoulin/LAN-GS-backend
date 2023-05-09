package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type EventsSheetService struct {
	spreadsheetId string
	readRange string
	srv *sheets.Service
}

type Event struct {
	Name string `mapstructure:"name"`	
	Description	string `mapstructure:"description"`
}

func NewEventsSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*EventsSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &EventsSheetService {
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (ESS *EventsSheetService) GetEvents(ctx context.Context) ([]Event, error) {
	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, ESS.readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("%v", err)
	}

	events := []Event{}

	colMap := map[int]string {
		0: "name",
		1: "description",
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