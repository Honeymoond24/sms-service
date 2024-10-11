package application

import "github.com/Honeymoond24/sms-service/internal/domain"

type ServicesRepository interface {
	GetServices() (map[string][]domain.Service, error)
	GetPhoneNumber(country, service string, sum int, exceptionPhoneSet []string) (string, int, error)
	//StoreSms(sms domain.SMS) error
}
