package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type B2bRequestSheetService struct {
	spreadsheetId string
	readRange     string
	srv           *sheets.Service
}

type B2bRequest struct {
	Id           string `json:"id" mapstructure:"id"`
	Date     	 string `json:"date" mapstructure:"date"`
	Name         string `json:"name" mapstructure:"name"`
	Email         string `json:"email" mapstructure:"email"`
	Phone         string `json:"phone" mapstructure:"phone"`
	Comment         string `json:"comment" mapstructure:"comment"`
}

type B2bRequestResponse struct {
	Id           string `json:"id" mapstructure:"id"`
	Date     	 string `json:"date" mapstructure:"date"`
	Name         string `json:"name" mapstructure:"name"`
	Email         string `json:"email" mapstructure:"email"`
	Phone         string `json:"phone" mapstructure:"phone"`
	Comment         string `json:"comment" mapstructure:"comment"`
}

func (e *B2bRequest) MarshalJSON() ([]byte, error) {
	aliasValue := struct {
		Id           string `json:"id" mapstructure:"id"`
		Date     	 string `json:"date" mapstructure:"date"`
		Name         string `json:"name" mapstructure:"name"`
		Email         string `json:"email" mapstructure:"email"`
		Phone         string `json:"phone" mapstructure:"phone"`
		Comment         string `json:"comment" mapstructure:"comment"`
	}{
		Id:           e.Id,
		Date:         e.Date,
		Name:         e.Name,
		Email:        e.Email,
		Phone:        e.Phone,
		Comment:      e.Comment,
	}
	return json.Marshal(aliasValue)
}

func NewB2bRequestSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*B2bRequestSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &B2bRequestSheetService{
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (ESS *B2bRequestSheetService) CreateB2bRequest(ctx context.Context, b2bRequest B2bRequest) error {
	readRange := "master!2:1000"

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
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
			b2bRequest.Name,
			b2bRequest.Email,
			b2bRequest.Phone,
			b2bRequest.Comment,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, newReadRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
