package application

import "github.com/Honeymoond24/sms-service/internal/domain"

type ServicesRepository interface {
	GetServices() (map[string]map[string]int, error)
	GetPhoneNumber(countryName, serviceName string, sum int, exceptionPhoneSet []string) (int, int, error)
	StoreSms(sms domain.SMS) error
	GetPhoneNumberByPhone(phone int) (domain.PhoneNumber, error)
	FinishActivation(activationId, status int) error
}
