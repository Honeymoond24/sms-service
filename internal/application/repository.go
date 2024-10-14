package application

import "github.com/Honeymoond24/sms-service/internal/domain"

type ServicesRepository interface {
	GetServices() (map[string][]domain.Service, error)
	GetPhoneNumber(countryName, serviceName string, sum int, exceptionPhoneSet []string) (string, int, error)
	StoreSms(sms domain.SMS) error
	GetPhoneNumberByPhone(phone int) (domain.PhoneNumber, error)
}
