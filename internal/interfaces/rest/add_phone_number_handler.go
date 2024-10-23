package rest

import (
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddPhoneNumber(c echo.Context, body map[string]interface{}, service *application.SmsService) error {
	c.Logger().Info("AddPhoneNumber", body)

	phones, ok := body["phones"].([]interface{})
	if !ok {
		return c.String(http.StatusBadRequest, "invalid phones list")
	}

	phonesList := make([]domain.PhoneNumber, 0)
	for _, phone := range phones {
		phoneMap, ok := phone.(map[string]interface{})
		if !ok {
			return c.String(http.StatusBadRequest, "invalid phone object")
		}
		countryName, ok := phoneMap["country"].(string)
		if !ok {
			return c.String(http.StatusBadRequest, "invalid country")
		}
		phoneNumber, ok := phoneMap["phone"].(float64)
		if !ok {
			return c.String(http.StatusBadRequest, "invalid phone")
		}
		phonesList = append(phonesList, domain.PhoneNumber{
			Country: countryName,
			Number:  int(phoneNumber),
		})
	}

	err := service.AddPhoneNumbers(phonesList)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "ERROR",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "SUCCESS",
	})
}
