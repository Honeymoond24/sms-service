package application

import (
	"errors"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockSmsRepository is a mock implementation of the SmsRepository interface
type MockSmsRepository struct {
	mock.Mock
}

func (m *MockSmsRepository) GetServices() (map[string]map[string]int, error) {
	args := m.Called()
	return args.Get(0).(map[string]map[string]int), args.Error(1)
}

func (m *MockSmsRepository) GetPhoneNumber(countryName, serviceName string, sum int, exceptionPhoneSet []string) (int, int, error) {
	args := m.Called(countryName, serviceName, sum, exceptionPhoneSet)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockSmsRepository) StoreSms(sms domain.SMS) error {
	args := m.Called(sms)
	return args.Error(0)
}

func (m *MockSmsRepository) GetPhoneNumberByPhone(phone int) (domain.PhoneNumber, error) {
	args := m.Called(phone)
	return args.Get(0).(domain.PhoneNumber), args.Error(1)
}

func (m *MockSmsRepository) FinishActivation(activationId, status int) error {
	args := m.Called(activationId, status)
	return args.Error(0)
}

func TestGetServices(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	mockRepo.On("GetServices").Return(map[string]map[string]int{"ru": {"vk": 1}}, nil)

	countries, err := service.GetServices()
	assert.NoError(t, err)
	assert.NotNil(t, countries)
	mockRepo.AssertExpectations(t)
}

func TestGetServices_Error(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	mockRepo.On("GetServices").Return(map[string]map[string]int{},
		errors.New("error while getting services"))

	countries, err := service.GetServices()
	assert.Error(t, err)
	assert.Nil(t, countries)
	mockRepo.AssertExpectations(t)
}

func TestGetNumber(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	mockRepo.On("GetPhoneNumber", "ru", "vk", 1, []string{"777"}).Return(77777777777, 123, nil)

	number, activationID, err := service.GetNumber("ru", "vk", 1, []string{"777"})
	assert.NoError(t, err)
	assert.NotZero(t, number)
	assert.NotZero(t, activationID)
	mockRepo.AssertExpectations(t)
}

func TestGetNumber_Error(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	mockRepo.On("GetPhoneNumber", "ru", "vk", 1, []string{"777"}).Return(
		77777777777, 123, errors.New("error while getting phone number"))

	number, activationID, err := service.GetNumber("ru", "vk", 1, []string{"777"})
	assert.Error(t, err)
	assert.Zero(t, number)
	assert.Zero(t, activationID)
	mockRepo.AssertExpectations(t)
}

func TestPushSms(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	sms := domain.SMS{
		ID:        1,
		PhoneTo:   domain.PhoneNumber{Number: 77777777777},
		PhoneFrom: "Sender",
		Text:      "Hello, World!",
	}

	mockRepo.On("GetPhoneNumberByPhone", 77777777777).Return(
		domain.PhoneNumber{Number: 77777777777}, nil)
	mockRepo.On("StoreSms", sms).Return(nil)

	err := service.PushSms(sms)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPushSms_Error(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	sms := domain.SMS{
		ID:        1,
		PhoneTo:   domain.PhoneNumber{Number: 77777777777},
		PhoneFrom: "Sender",
		Text:      "Hello, World!",
	}

	mockRepo.On("GetPhoneNumberByPhone", 77777777777).Return(
		domain.PhoneNumber{Number: 77777777777}, nil)
	mockRepo.On("StoreSms", sms).Return(errors.New("error while getting phone number"))

	err := service.PushSms(sms)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFinishActivation(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	activationId := 1
	status := 3

	mockRepo.On("FinishActivation", activationId, status).Return(nil)

	err := service.FinishActivation(activationId, status)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFinishActivation_Error(t *testing.T) {
	mockRepo := new(MockSmsRepository)
	service := NewSmsService(mockRepo)

	activationId := 1
	status := 3

	mockRepo.On("FinishActivation", activationId, status).Return(errors.New("update error"))

	err := service.FinishActivation(activationId, status)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
