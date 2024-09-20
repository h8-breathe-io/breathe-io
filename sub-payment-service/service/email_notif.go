package service

type EmailNotifService interface {
	NotifyPaymentSucccess()
}

func NewEmailNotifService() EmailNotifService {
	return &emailNotifService{}
}

type emailNotifService struct {
}

func (es *emailNotifService) NotifyPaymentSucccess() {

}
