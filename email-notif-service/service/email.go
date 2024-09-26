package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SendEmailRequest struct {
	From     Address   `json:"from"`
	To       []Address `json:"to"`
	Subject  string    `json:"subject"`
	Text     string    `json:"text"`
	Html     string    `json:"html"`
	Category string    `json:"category"`
}

type EmailService interface {
	SendEmail(req *SendEmailRequest) error
}

type emailService struct {
	client *http.Client
	url    string
}

func NewEmailService(url string) EmailService {
	return &emailService{
		client: &http.Client{},
		url:    url,
	}
}

func (es *emailService) SendEmail(req *SendEmailRequest) error {
	method := "POST"
	encoded, err := json.Marshal(req)
	payload := bytes.NewReader(encoded)
	if err != nil {

		return status.Errorf(codes.Internal, "error encoding mail request: %s", err.Error())
	}
	httpreq, err := http.NewRequest(method, es.url, payload)
	if err != nil {
		return status.Errorf(codes.Internal, "error creating mail request: %s", err.Error())
	}
	httpreq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("MAILTRAP_TOKEN")))
	httpreq.Header.Add("Content-Type", "application/json")

	res, err := es.client.Do(httpreq)

	if err != nil {

		return status.Errorf(codes.Internal, "error calling mail api: %s", err.Error())
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return status.Errorf(codes.Internal, "error reading mail api response: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return status.Errorf(codes.Internal, "error mail api request: %s", string(resp))
	}
	log.Print("mail api success: %s", string(resp))
	return nil
}
