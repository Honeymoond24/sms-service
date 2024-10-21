package application

import (
	"errors"
	"fmt"
	"github.com/Honeymoond24/sms-service/internal/config"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"log"
	"math/rand"
	"time"
)

type SmsService struct {
	repo ServicesRepository
	cfg  *config.Config
}

func NewSmsService(servicesRepo ServicesRepository, cfg *config.Config) *SmsService {
	return &SmsService{servicesRepo, cfg}
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

	go func() {
		sms := domain.SMS{
			ID:        rand.Intn(1023),
			PhoneTo:   domain.PhoneNumber{Number: number},
			PhoneFrom: serviceName,
			Text:      fmt.Sprintf("%s activation code: %d", serviceName, time.Now().Nanosecond()),
		}
		err := s.PushSms(sms)
		if err != nil {
			log.Println(err)
		}
	}()

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
