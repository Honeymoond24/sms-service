package application

import (
	"errors"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"log"
)

type SmsService struct {
	repo ServicesRepository
}

func NewSmsService(servicesRepo ServicesRepository) *SmsService {
	return &SmsService{servicesRepo}
}

func (s *SmsService) GetServices() (map[string]map[string]int, error) {
	countries, err := s.repo.GetServices()
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while getting services")
	}

	return countries, nil
}

func (s *SmsService) GetNumber(countryName, serviceName string, sum int, phonePrefixes []string) (int, int, error) {
	number, activationID, err := s.repo.GetPhoneNumber(
		countryName,
		serviceName,
		sum,
		phonePrefixes,
	)
	if err != nil {
		if errors.Is(err, PhoneNotFound) {
			return 0, 0, PhoneNotFound
		}
		return 0, 0, errors.New("error while getting phone number")
	}

	return number, activationID, nil
}

func (s *SmsService) PushSms(sms domain.SMS) error {
	phoneNumber, err := s.repo.GetPhoneNumberByPhone(sms.PhoneTo.Number)
	if err != nil {
		log.Println(err)
		return errors.New("error while getting phone number")
	}
	if phoneNumber.Number == 0 {
		return PhoneNotFound
	}
	sms.PhoneTo = phoneNumber
	err = s.repo.StoreSms(sms)
	if err != nil {
		log.Println(err)
		return errors.New("error while storing sms")
	}

	return nil
}

func (s *SmsService) FinishActivation(activationId, status int) error {
	return s.repo.FinishActivation(activationId, status)
}
