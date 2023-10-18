package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	ColID = 0
	ColName = 1
	ColDate = 2
	ColDescription = 3
	ColLink = 4
	ColExternalLink = 5
	ColCapacity = 6
	ColType = 7
	ColShowForm = 8 
)

type EventsSheetService struct {
	spreadsheetId string
	readRange     string
	srv           *sheets.Service
}

type Event struct {
	Id           string `json:"id" mapstructure:"id"`
	Name         string `json:"name" mapstructure:"name"`
	Date     	 string `json:"date" mapstructure:"date"`
	Description  string `json:"description" mapstructure:"description"`
	Link         string `json:"link" mapstructure:"link"`
	ExternalLink string `json:"externalLink" mapstructure:"externalLink"`
	Capacity     string `json:"capacity" mapstructure:"capacity"`
	Type     	 string `json:"type" mapstructure:"type"`
	ShowForm     bool `json:"showForm" mapstructure:"showForm"`
}

type EventResponse struct {
	Id           string `json:"id" mapstructure:"id"`
	Name         string `json:"name" mapstructure:"name"`
	Date     	 string `json:"date" mapstructure:"date"`
	Description  string `json:"description" mapstructure:"description"`
	Link         string `json:"link" mapstructure:"link"`
	ExternalLink string `json:"externalLink" mapstructure:"externalLink"`
	Capacity     string `json:"capacity" mapstructure:"capacity"`
	Type     	 string `json:"type" mapstructure:"type"`
	ShowForm     bool `json:"showForm" mapstructure:"showForm"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
	aliasValue := struct {
		Id           string `json:"id" mapstructure:"id"`
		Name         string `json:"name" mapstructure:"name"`
		Date         string `json:"date" mapstructure:"date"`
		Description  string `json:"description" mapstructure:"description"`
		Link         string `json:"link" mapstructure:"link"`
		ExternalLink string `json:"externalLink" mapstructure:"externalLink"`
		Capacity     string `json:"capacity" mapstructure:"capacity"`
		Type     	 string `json:"type" mapstructure:"type"`
		ShowForm     bool   `json:"showForm" mapstructure:"showForm"`
	}{
		Id:           e.Id,
		Name:         e.Name,
		Date:         e.Date,
		Description:  e.Description,
		Link:         e.Link,
		ExternalLink: e.ExternalLink,
		Capacity:     e.Capacity,
		Type: 		  e.Type,
		ShowForm:     e.ShowForm,
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
		Type:     	  ev.Type,
		ShowForm:     ev.ShowForm,
	}

	return nil
}

func isNotExpired(date string) bool {
	d, _ := time.Parse("02.01.2006", date)
	now := time.Now()

	return d.Add(time.Hour * 24).After(now)
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
		2: "date",
		3: "description",
		4: "link",
		5: "externalLink",
		6: "capacity",
		7: "type",
		8: "showForm",
	}

	for _, val := range res.Values {
		e := map[string]interface{}{}

		showForm, _ := strconv.ParseBool(val[ColShowForm].(string))

		for i, v := range val {
			col, ok := colMap[i]
			if ok {
				e[col] = v
			}

			if i == ColShowForm {
				e[col] = showForm;
			}
		}

		var event Event

		mapstructure.Decode(e, &event)

		if (isNotExpired(event.Date)) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (ESS *EventsSheetService) GetEvent(ctx context.Context, eventId string) (Event, error) {
	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, ESS.readRange).Do()

	event := Event{}

	if err != nil || res.HTTPStatusCode != 200 {
		return event, fmt.Errorf("%v", err)
	}


	colMap := map[int]string{
		0: "id",
		1: "name",
		2: "date",
		3: "description",
		4: "link",
		5: "externalLink",
		6: "capacity",
		7: "type",
		8: "showForm",
	}

	for _, val := range res.Values {
		e := map[string]interface{}{}

		for i, v := range val {
			col, ok := colMap[i]
			if ok {
				e[col] = v
			}
		}

		var eventTmp Event

		mapstructure.Decode(e, &eventTmp)

		if (isNotExpired(eventTmp.Date) && eventTmp.Id == eventId) {
			event = eventTmp
		}
	}

	return event, nil
}

func (ESS *EventsSheetService) UpdateEvent(ctx context.Context, event Event) error {
	var rowNumber int
	readRange := "master!2:1000"

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
	}

	for i, v := range res.Values {
		if len(v) > 0 && v[0] == event.Id {
			rowNumber = i + 2
		}
	}

	updateRowRange := fmt.Sprintf("A%d:I%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			event.Id,
			event.Name,
			event.Date,
			event.Description,
			event.Link,
			event.ExternalLink,
			event.Capacity,
			event.Type,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
