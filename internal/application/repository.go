package application

import "github.com/Honeymoond24/sms-service/internal/domain"

type ServicesRepository interface {
	CreateService(service domain.Service) (int, error)
	GetServices() ([]domain.Service, error)
	DeleteService(id int) error
}

type PhoneNumberRepository interface {
	CreatePhoneNumber(phoneNumber domain.PhoneNumber) (int, error)
}

type SmsRepository interface {
	StoreSms(sms domain.SMS) error
}
