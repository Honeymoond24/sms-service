package rest

import (
	"errors"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PushSms(c echo.Context, body map[string]interface{}, service *application.SmsService) error {
	c.Logger().Info("PushSms", body)

	smsIdFloat, ok := body["smsID"].(float64)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid smsID")
	}
	smsId := int(smsIdFloat)

	phoneFloat, ok := body["phone"].(float64)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid phone")
	}
	phone := int(phoneFloat)

	phoneFrom, ok := body["phoneFrom"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid phoneFrom")
	}

	text, ok := body["text"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid text")
	}

	sms := domain.SMS{
		ID:        smsId,
		Phone:     phone,
		PhoneFrom: phoneFrom,
		Text:      text,
	}

	err := service.PushSms(sms)
	if err != nil {
		response := map[string]interface{}{"status": "ERROR"}
		if errors.Is(err, application.PhoneNotFound) {
			response["error"] = err.Error()
			return c.JSON(http.StatusNotFound, response)
		}
		response["error"] = "internal server error"
		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "SUCCESS",
	})
}
