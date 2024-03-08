package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	CreationDate string `json:"creationDate" mapstructure:"creationDate"`
	Name string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
	Phone string `json:"phone" mapstructure:"phone"`
	NumberOfPersons string `json:"numberOfPersons" mapstructure:"numberOfPersons"`
	Instagram string `json:"instagram" mapstructure:"instagram"`
	Telegram string `json:"telegram" mapstructure:"telegram"`
	Date string `json:"date" mapstructure:"date"`
	EventId string `json:"eventId" mapstructure:"eventId"`
	Comment	string `json:"comment" mapstructure:"comment"`
	WillCome bool `json:"willCome" mapstructure:"willCome"`
}

type EntrieResponse struct {
	Id string `json:"id" mapstructure:"id"`
}

type CancelEntrieResponse struct {
	Id string `json:"id" mapstructure:"id"`
}

const (
	Id = 0
	CreationDate = 1
	Name = 2
	Email = 3
	Phone = 4
	NumberOfPersons = 5
	Instagram = 6
	Telegram = 7
	Date = 8
	EventId = 9
	Comment = 10
	WillCome = 11
)

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

	newReadRange := fmt.Sprintf("A%v:L%v", tableLen, tableLen)

	id := uuid.New()

	now := time.Now()
	formatted := now.Format(time.RFC3339)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			id.String(),
			formatted,
			entrie.Name,
			entrie.Email,
			entrie.Phone,
			entrie.NumberOfPersons,
			entrie.Instagram,
			entrie.Telegram,
			entrie.Date,
			entrie.EventId,
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
		if v[Id] == entrie.Id {
			rowNumber = i + 2
		}
	}

	updateRowRange := fmt.Sprintf("A%d:L%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			entrie.Id,
			entrie.CreationDate,
			entrie.Name,
			entrie.Email,
			entrie.Phone,
			entrie.NumberOfPersons,
			entrie.Instagram,
			entrie.Telegram,
			entrie.Date,
			entrie.EventId,
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
		fmt.Println(v[EventId], cancelEntrie.Id)
		if v[EventId] == cancelEntrie.Id {
			rowNumber = i + 2
			entrie = v
		}
	}

	if entrie == nil {
		return nil
	}

	updateRowRange := fmt.Sprintf("A%d:L%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			entrie[Id],
			entrie[CreationDate],
			entrie[Name],
			entrie[Email],
			entrie[Phone],
			entrie[NumberOfPersons],
			entrie[Instagram],
			entrie[Telegram],
			entrie[Date],
			entrie[EventId],
			entrie[Comment],
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
		if len(v) >= 11 {
			eventId, ok := v[EventId].(string)
		
			if ok {
				willCome, ok := v[WillCome].(string)
	
				if ok {
					willCome, _ := strconv.ParseBool(willCome)
	
					if willCome {
						_, ok := uniqueEntries[eventId]
			
						outInt, _ := strconv.Atoi(v[NumberOfPersons].(string))
				
						if ok {
							uniqueEntries[eventId] += outInt
						} else {
							uniqueEntries[eventId] = outInt
						}
					}
				}
			}
		}
	}

	return uniqueEntries, nil
}
