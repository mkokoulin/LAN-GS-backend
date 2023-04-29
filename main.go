package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Secret struct {
	Type string `json:"type"`
	ProjectId string `json:"project_id"`
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	ClientId string `json:"client_id"`
	AuthUri string `json:"auth_uri"`
	TokenUri string `json:"token_uri"`
}

func main() {
	// create api context
	ctx := context.Background()

	data, err := os.ReadFile("lan-site-94255-b56432cf737d.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

	fmt.Println("1")

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("2")

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("3")

	res, err := srv.Spreadsheets.Get("1DbNF1xd8Un3GwCWIhTsFfsdazIZC3kSOLOSw0byulJA").Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || res.HTTPStatusCode != 200 {
		log.Fatal(err)
		return
	}
	fmt.Println("4")

	for _, v := range res.Sheets {
		prop := v.Properties
		fmt.Println(prop.Title)
	}
}