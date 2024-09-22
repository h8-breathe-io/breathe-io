package server

import (
	"context"
	"email-notif-service/pb"
	"email-notif-service/service"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewEmailNotifServer(
	paymentService service.SubsPaymentService,
	userService service.UserService,
) *EmailNotifServer {
	return &EmailNotifServer{
		paymentService: paymentService,
		userService:    userService,
	}
}

type EmailNotifServer struct {
	paymentService service.SubsPaymentService
	userService    service.UserService
	pb.UnimplementedEmailNotifServiceServer
}

func (es *EmailNotifServer) NotifyPaymentComplete(c context.Context, req *pb.NotifyPaymentCompleteReq) (*pb.NotifyPaymentCompleteResp, error) {
	// TODO
	log.Printf("received notify payment complete request")
	url := "https://sandbox.api.mailtrap.io/api/send/3155742"
	method := "POST"
	payload := strings.NewReader(`{"from":{"email":"hello@example.com","name":"Mailtrap Test"},"to":[{"email":"razif.dev@gmail.com"}],"subject":"You are awesome!","text":"Congrats for sending test email with Mailtrap!","category":"Integration Test"}`)

	client := &http.Client{}
	httpreq, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "error calling mail api: %s", err.Error())
	}
	httpreq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("MAILTRAP_TOKEN")))
	httpreq.Header.Add("Content-Type", "application/json")

	res, err := client.Do(httpreq)

	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "error calling mail api: %s", err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "error calling mail api: %s", err.Error())
	}

	fmt.Println(string(body))
	return &pb.NotifyPaymentCompleteResp{
		Status: "OK",
		Email:  "razif.dev@gmail.com",
	}, nil
}
