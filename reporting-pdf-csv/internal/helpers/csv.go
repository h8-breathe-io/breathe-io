package helpers

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// WriteToGoogleSheet writes the provided data to Google Sheets and returns the sheet URL
func WriteToGoogleSheet(sheetTitle string, data [][]string) (string, error) {
	// Load credentials from credentials.json
	creds, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return "", fmt.Errorf("unable to read client secret file: %v", err)
	}

	// Authenticate using service account credentials
	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope, drive.DriveFileScope)
	if err != nil {
		return "", fmt.Errorf("failed to create config from JSON: %v", err)
	}

	client := config.Client(context.Background())

	// Create Google Sheets and Drive services
	srvSheets, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	srvDrive, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	// Create a new Google Sheet
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: sheetTitle,
		},
	}

	createResp, err := srvSheets.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create spreadsheet: %v", err)
	}

	spreadsheetID := createResp.SpreadsheetId
	sheetURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID)

	// Convert [][]string to [][]interface{} (as required by Sheets API)
	var vr sheets.ValueRange
	for _, row := range data {
		convertedRow := make([]interface{}, len(row))
		for i, cell := range row {
			convertedRow[i] = cell
		}
		vr.Values = append(vr.Values, convertedRow)
	}

	// Write data to the Google Sheet
	_, err = srvSheets.Spreadsheets.Values.Update(spreadsheetID, "Sheet1", &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return "", fmt.Errorf("unable to write data to sheet: %v", err)
	}

	// Make the Google Sheet accessible to anyone with the link in read-only mode
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err = srvDrive.Permissions.Create(spreadsheetID, permission).Do()
	if err != nil {
		return "", fmt.Errorf("unable to make sheet public: %v", err)
	}

	return sheetURL, nil
}
