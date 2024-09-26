package handlers

import (
	"context"
	"fmt"
	"reporting/internal/helpers"
	"reporting/internal/service"
	"reporting/proto/pb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ReportService struct holds necessary dependencies, such as database connection
type ReportService struct {
	DB                                  *gorm.DB
	UserService                         service.UserService
	pb.UnimplementedReportServiceServer // This is required to avoid breaking gRPC compatibility
}

// GenerateReport handles the request for generating both CSV and PDF reports and returns their links
func (s *ReportService) GenerateReport(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	var result []struct {
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
	}

	// Parsing start_date and end_date if provided
	var startDate, endDate time.Time
	var err error

	// validate token and get user
	user, err := s.UserService.ValidateAndGetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid token '%s'", err.Error())
	}

	if req.StartDate != "" {
		startDate, err = time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format: %v", err)
		}
	}
	if req.EndDate != "" {
		endDate, err = time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %v", err)
		}
	}

	userID := user.ID

	// Build the query
	query := s.DB.Table("users u").
		Select("u.username, u.email, u.phonenumber, u.tier, "+
			"COALESCE(SUM(bf.total_emission), 0) AS total_business_emission, "+
			"bf.company_type, bf.location_id AS business_location_id, bf.total_emission, "+
			"l.location_name, aq.aqi, aq.co, aq.no, aq.no2, aq.o3, aq.so2, aq.pm25, aq.pm10, aq.nh3").
		Joins("LEFT JOIN business_facilities bf ON u.id = bf.user_id").
		Joins("LEFT JOIN user_location ul ON u.id = ul.user_id").
		Joins("LEFT JOIN locations l ON ul.location_id = l.id").
		Joins("LEFT JOIN air_quality aq ON l.id = aq.location_id").
		Where("u.id = ?", userID)

	// Add date filtering if start_date and end_date are provided
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("aq.fetch_time BETWEEN ? AND ?", startDate, endDate)
	}

	// Execute the query
	if err := query.Group("u.id, bf.company_type, bf.location_id, bf.total_emission, l.location_name, aq.aqi, aq.co, aq.no, aq.no2, aq.o3, aq.so2, aq.pm25, aq.pm10, aq.nh3").Scan(&result).Error; err != nil {
		return nil, err
	}

	// Prepare data for CSV and PDF
	sheetData := [][]string{
		{"Username", "Email", "Phone Number", "Tier", "Total Business Emission", "Company Type", "Business Location ID", "Total Emission", "Location Name", "AQI", "CO", "NO", "NO2", "O3", "SO2", "PM25", "PM10", "NH3"},
	}

	for _, r := range result {
		sheetData = append(sheetData, []string{
			r.Username,
			r.Email,
			r.PhoneNumber,
			r.Tier,
			fmt.Sprintf("%.2f", r.TotalBusinessEmission),
			r.CompanyType,
			fmt.Sprintf("%d", r.BusinessLocationID),
			fmt.Sprintf("%.2f", r.TotalEmission),
			r.LocationName,
			fmt.Sprintf("%d", r.AQI),
			fmt.Sprintf("%.2f", r.CO),
			fmt.Sprintf("%.2f", r.NO),
			fmt.Sprintf("%.2f", r.NO2),
			fmt.Sprintf("%.2f", r.O3),
			fmt.Sprintf("%.2f", r.SO2),
			fmt.Sprintf("%.2f", r.PM25),
			fmt.Sprintf("%.2f", r.PM10),
			fmt.Sprintf("%.2f", r.NH3),
		})
	}

	// Generate human-readable timestamp for the filenames
	timestamp := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	csvFileName := fmt.Sprintf("report_%s_%s", req.Tier, timestamp)
	pdfFileName := fmt.Sprintf("report_%s_%s.pdf", req.Tier, timestamp)

	// Step 1: Generate and Upload CSV to Google Sheets
	csvURL, err := helpers.WriteToGoogleSheet(csvFileName, sheetData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate or upload CSV: %v", err)
	}

	// Step 2: Generate Prompt for ChatGPT based on query data
	prompt := helpers.BuildPromptFromData(result)

	// Step 3: Send the prompt to ChatGPT and retrieve recommendations
	suggestions, err := helpers.GetSuggestionsFromChatGPT(ctx, prompt) // Use `suggestions` instead of `gptResponse`
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions from ChatGPT: %v", err)
	}

	// Step 4: Generate the PDF combining query results and ChatGPT suggestions
	pdfContent, err := helpers.GeneratePDF(result, []string{suggestions}) // Pass `suggestions` as a slice of strings
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	// Step 5: Upload the PDF to Google Drive
	pdfURL, err := helpers.UploadToGoogleDrive(pdfContent, pdfFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload PDF to Google Drive: %v", err)
	}

	// Return the URLs for both CSV and PDF
	return &pb.ReportResponse{
		CsvLocation: csvURL,
		PdfLocation: pdfURL,
	}, nil
}
