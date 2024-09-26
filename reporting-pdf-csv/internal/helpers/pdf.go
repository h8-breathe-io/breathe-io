package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// BuildPromptFromData converts query result data into a prompt for ChatGPT
func BuildPromptFromData(result []struct {
	Username              string
	Email                 string
	PhoneNumber           string
	Tier                  string
	TotalBusinessEmission float64
	CompanyType           string
	BusinessLocationID    uint
	TotalEmission         float64
	LocationName          string
	AQI                   int
	CO                    float64
	NO                    float64
	NO2                   float64
	O3                    float64
	SO2                   float64
	PM25                  float64
	PM10                  float64
	NH3                   float64
}) string {
	// Build the prompt for ChatGPT based on query results
	prompt := "Please provide suggestions and recommendations based on the following data:\n\n"
	for _, r := range result {
		prompt += fmt.Sprintf("User: %s, Tier: %s, Total Business Emission: %.2f, AQI: %d, Location: %s\n",
			r.Username, r.Tier, r.TotalBusinessEmission, r.AQI, r.LocationName)
	}
	prompt += "\nProvide suggestions for improvement in air quality and business emissions management."
	return prompt
}

// GetSuggestionsFromChatGPT sends the prompt to ChatGPT and retrieves recommendations using CreateChatCompletion
func GetSuggestionsFromChatGPT(ctx context.Context, prompt string) (string, error) {
	client := openai.NewClient(os.Getenv("TOKEN_GPT")) // Ganti dengan API key Anda

	// Create a chat-based request
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo", // You can also use "gpt-3.5-turbo"
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system", // Optional: You can give the model some instructions about how it should behave.
				Content: "You are an assistant that provides helpful suggestions.",
			},
			{
				Role:    "user", // This is the actual user prompt
				Content: prompt,
			},
		},
	}

	// Get the chat-based response
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error getting suggestions from ChatGPT: %v", err)
	}

	// Return the first response from the assistant
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no suggestions returned by ChatGPT")
}

// GeneratePDF generates a PDF with the provided data and suggestions from ChatGPT
func GeneratePDF(result []struct {
	Username              string
	Email                 string
	PhoneNumber           string
	Tier                  string
	TotalBusinessEmission float64
	CompanyType           string
	BusinessLocationID    uint
	TotalEmission         float64
	LocationName          string
	AQI                   int
	CO                    float64
	NO                    float64
	NO2                   float64
	O3                    float64
	SO2                   float64
	PM25                  float64
	PM10                  float64
	NH3                   float64
}, suggestions []string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(40, 10, "User Report")
	pdf.Ln(12)

	// Add query data to the PDF
	pdf.SetFont("Arial", "", 12)
	for _, r := range result {
		pdf.Cell(0, 10, fmt.Sprintf("Username: %s, Email: %s, Tier: %s", r.Username, r.Email, r.Tier))
		pdf.Ln(10)
		pdf.Cell(0, 10, fmt.Sprintf("Total Business Emission: %.2f, AQI: %d", r.TotalBusinessEmission, r.AQI))
		pdf.Ln(10)
		pdf.Cell(0, 10, fmt.Sprintf("Location: %s", r.LocationName))
		pdf.Ln(10)
	}

	// Add ChatGPT suggestions and recommendations
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Suggestions and Recommendations")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for _, suggestion := range suggestions {
		pdf.MultiCell(0, 10, suggestion, "", "L", false)
		pdf.Ln(8)
	}

	// Save the PDF to a byte buffer
	var buffer bytes.Buffer
	err := pdf.Output(&buffer)
	if err != nil {
		return nil, fmt.Errorf("error generating PDF: %v", err)
	}

	return buffer.Bytes(), nil
}

// UploadToGoogleDrive uploads the provided PDF file content to Google Drive and returns the public URL
func UploadToGoogleDrive(fileData []byte, fileName string) (string, error) {
	creds, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return "", err
	}

	config, err := google.JWTConfigFromJSON(creds, drive.DriveFileScope)
	if err != nil {
		return "", err
	}

	client := config.Client(context.Background())
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return "", err
	}

	file := &drive.File{Name: fileName, MimeType: "application/pdf"}
	uploadFile := bytes.NewReader(fileData)
	uploadedFile, err := srv.Files.Create(file).Media(uploadFile).Do()
	if err != nil {
		return "", err
	}

	permission := &drive.Permission{Type: "anyone", Role: "reader"}
	_, err = srv.Permissions.Create(uploadedFile.Id, permission).Do()
	if err != nil {
		return "", err
	}

	return "https://drive.google.com/file/d/" + uploadedFile.Id + "/view", nil
}
