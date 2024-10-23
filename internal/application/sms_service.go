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

func (s *SmsService) GetNumber(
	countryName, serviceName string, sum int, phonePrefixes []string,
	callback func(url string, sms domain.SMS),
) (int, int, error) {
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
		// simulate sms sending process with random delay from 0 to 10 seconds
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

		sms := domain.SMS{
			ID:        rand.Intn(1023),
			PhoneTo:   domain.PhoneNumber{Number: number},
			PhoneFrom: serviceName,
			Text:      fmt.Sprintf("%s activation code: %d", serviceName, time.Now().Nanosecond()),
		}
		log.Println("sms", sms)
		err := s.PushSms(sms, callback)
		if err != nil {
			log.Println(err)
		}
	}()

	return number, activationID, nil
}

func (s *SmsService) PushSms(sms domain.SMS, callback func(url string, sms domain.SMS)) error {
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

	callback(s.cfg.PushSMSURL, sms)

	return nil
}

func (s *SmsService) FinishActivation(activationId, status int) error {
	return s.repo.FinishActivation(activationId, status)
}

func (s *SmsService) AddPhoneNumbers(phoneNumbers []domain.PhoneNumber) error {
	return s.repo.AddPhoneNumbers(phoneNumbers)
}
