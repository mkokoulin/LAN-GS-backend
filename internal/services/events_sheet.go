package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type EventsSheetService struct {
	spreadsheetId string
	readRange     string
	srv           *sheets.Service
}

type Event struct {
	Id           string `json:"id" mapstructure:"id"`
	Name         string `json:"name" mapstructure:"name"`
	Description  string `json:"description" mapstructure:"description"`
	Link         string `json:"link" mapstructure:"link"`
	ExternalLink string `json:"externalLink" mapstructure:"externalLink"`
	Capacity     string `json:"capacity" mapstructure:"capacity"`
}

type EventResponse struct {
	Id           string `json:"id" mapstructure:"id"`
	Name         string `json:"name" mapstructure:"name"`
	Description  string `json:"description" mapstructure:"description"`
	Link         string `json:"link" mapstructure:"link"`
	ExternalLink string `json:"externalLink" mapstructure:"externalLink"`
	Capacity     string `json:"capacity" mapstructure:"capacity"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
	aliasValue := struct {
		Id           string `json:"id" mapstructure:"id"`
		Name         string `json:"name" mapstructure:"name"`
		Description  string `json:"description" mapstructure:"description"`
		Link         string `json:"link" mapstructure:"link"`
		ExternalLink string `json:"externalLink" mapstructure:"externalLink"`
		Capacity     string `json:"capacity" mapstructure:"capacity"`
	}{
		Id:           e.Id,
		Name:         e.Name,
		Description:  e.Description,
		Link:         e.Link,
		ExternalLink: e.ExternalLink,
		Capacity:     e.Capacity,
	}
	return json.Marshal(aliasValue)
}

func (e *Event) UnmarshalJSON(b []byte) error {
	var ev EventResponse

	if err := json.Unmarshal(b, &ev); err != nil {
		return err
	}

	*e = Event{
		Id:           ev.Id,
		Name:         ev.Name,
		Description:  ev.Description,
		Link:         ev.Link,
		ExternalLink: ev.ExternalLink,
		Capacity:     ev.Capacity,
	}

	return nil
}

func NewEventsSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*EventsSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &EventsSheetService{
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

	colMap := map[int]string{
		0: "id",
		1: "name",
		2: "description",
		3: "link",
		4: "externalLink",
		5: "capacity",
	}

	for _, val := range res.Values {
		e := map[string]interface{}{}

		for i, v := range val {
			col, ok := colMap[i]
			if ok {
				e[col] = v
			}
		}

		var event Event

		mapstructure.Decode(e, &event)

		events = append(events, event)
	}

	return events, nil
}

func (ESS *EventsSheetService) UpdateEvent(ctx context.Context, event Event) error {
	var rowNumber int
	readRange := "master!2:1000"

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
	}

	for i, v := range res.Values {
		if v[0].(string)+v[1].(string) == event.Name+event.Description {
			rowNumber = i + 2
		}
	}

	updateRowRange := fmt.Sprintf("A%d:E%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			event.Id,
			event.Name,
			event.Description,
			event.Link,
			event.ExternalLink,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
