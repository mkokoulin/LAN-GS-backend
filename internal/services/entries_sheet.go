package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type EntriesSheetService struct {
	spreadsheetId string
	srv *sheets.Service
}

type Entrie struct {
	Id string `json:"id" mapstructure:"id"`
	Name string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
	Phone string `json:"phone" mapstructure:"phone"`
	NumberOfPersons string `json:"numberOfPersons" mapstructure:"numberOfPersons"`
	Instagram string `json:"instagram" mapstructure:"instagram"`
	Telegram string `json:"telegram" mapstructure:"telegram"`
	Date string `json:"date" mapstructure:"date"`
	Event string `json:"event" mapstructure:"event"`
	Comment	string `json:"comment" mapstructure:"comment"`
	WillCome bool `json:"willCome" mapstructure:"willCome"`
}

type EntrieResponse struct {
	Id string `json:"id" mapstructure:"id"`
}

type CancelEntrieResponse struct {
	Id string `json:"id" mapstructure:"id"`
}

func NewEntriesSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*EntriesSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &EntriesSheetService {
		spreadsheetId,
		srv,
	}, nil
}

func (ESS *EntriesSheetService) CreateEntrie(ctx context.Context, entrie Entrie) (EntrieResponse, error) {
	readRange := "master!2:1000"
	response := EntrieResponse {}

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return response, fmt.Errorf("%v", err)
	}

	tableLen := len(res.Values) + 2

	newReadRange := fmt.Sprintf("A%v:J%v", tableLen, tableLen)

	id := uuid.New()

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			id.String(),
			entrie.Name,
			entrie.Email,
			entrie.Phone,
			entrie.NumberOfPersons,
			entrie.Instagram,
			entrie.Telegram,
			entrie.Date,
			entrie.Event,
			entrie.Comment,
			true,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, newReadRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return response, fmt.Errorf("%v", err)
	}

	response.Id = id.String()

	return response, nil
}

func (ESS *EntriesSheetService) UpdateEntrie(ctx context.Context, entrie Entrie) error {
	var rowNumber int
	readRange := "master!2:1000" 

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
	}

	for i, v := range res.Values {
		if v[0] == entrie.Id {
			rowNumber = i + 2
		}
	}

	updateRowRange := fmt.Sprintf("A%d:J%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			entrie.Id,
			entrie.Name,
			entrie.Email,
			entrie.Phone,
			entrie.NumberOfPersons,
			entrie.Instagram,
			entrie.Telegram,
			entrie.Date,
			entrie.Event,
			entrie.Comment,
			entrie.WillCome,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (ESS *EntriesSheetService) CancelEntrie(ctx context.Context, cancelEntrie CancelEntrieResponse) error {
	var rowNumber int
	var entrie []interface{}
	readRange := "master!2:1000" 

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
	}

	for i, v := range res.Values {
		if v[7] == cancelEntrie.Id {
			rowNumber = i + 2
			entrie = v
		}
	}

	if entrie == nil {
		return nil
	}

	updateRowRange := fmt.Sprintf("A%d:J%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			entrie[0],
			entrie[1],
			entrie[2],
			entrie[3],
			entrie[4],
			entrie[5],
			entrie[6],
			entrie[7],
			entrie[8],
			false,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (ESS *EntriesSheetService) GetUniqueEntries(ctx context.Context) (map[string]int, error) {
	uniqueEntries := map[string]int {}
	readRange := "master!2:1000"

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("%v", err)
	}

	for _, v := range res.Values {
		eventId := v[8]

		_, ok := uniqueEntries[eventId.(string)]

		if ok {
			uniqueEntries[eventId.(string)] += 1
		} else {
			uniqueEntries[eventId.(string)] = 1
		}
	}

	return uniqueEntries, nil
}