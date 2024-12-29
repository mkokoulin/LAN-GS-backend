package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type BookingSheetService struct {
	spreadsheetId string
	srv *sheets.Service
}

type Booking struct {
	Id string `json:"id" mapstructure:"id"`
	CreatedAt string `json:"createdAt" mapstructure:"createdAt"`
	Name string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
	Phone string `json:"phone" mapstructure:"phone"`
	NumberOfPersons string `json:"numberOfPersons" mapstructure:"numberOfPersons"`
	Instagram string `json:"instagram" mapstructure:"instagram"`
	Telegram string `json:"telegram" mapstructure:"telegram"`
	Date string `json:"date" mapstructure:"date"`
	Comment	string `json:"comment" mapstructure:"comment"`
}

func NewBookingSheets(ctx context.Context, googleClient *http.Client, spreadsheetId string) (*BookingSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &BookingSheetService {
		spreadsheetId,
		srv,
	}, nil
}

func (ESS *BookingSheetService) CreateBooking(ctx context.Context, booking Booking) error {
	readRange := "master!2:10000"

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
			booking.Name,
			booking.Email,
			booking.Phone,
			booking.NumberOfPersons,
			booking.Instagram,
			booking.Telegram,
			booking.Date,
			booking.Comment,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, newReadRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}