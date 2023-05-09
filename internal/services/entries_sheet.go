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
	Id string `mapstructure:"id"`
	Name string `mapstructure:"name"`
	Email string `mapstructure:"email"`
	Phone string `mapstructure:"phone"`
	NumberOfPersons string `mapstructure:"numberOfPersons"`
	Social string `mapstructure:"social"`
	Date string `mapstructure:"date"`
	Event string `mapstructure:"event"`
	Comment	string `mapstructure:"comment"`
	WillCome string `mapstructure:"willCome"`
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

func (ESS *EntriesSheetService) CreateEntrie(ctx context.Context, entrie Entrie) error {
	readRange := "master!2:1000" 

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
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
			entrie.Social,
			entrie.Date,
			entrie.Event,
			entrie.Comment,
			true,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, newReadRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}