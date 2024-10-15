package rest

import (
	"errors"
	"fmt"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"net/http"
)

func FinishActivation(c echo.Context, body map[string]interface{}, service *application.SmsService) error {
	c.Logger().Info("FinishActivation", body)

	activationIdFloat, ok := body["activationId"].(float64)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid activationId")
	}
	activationId := int(activationIdFloat)

	statusFloat, ok := body["status"].(float64)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid status")
	}
	status := int(statusFloat)
	fmt.Println(activationId, status)

	err := service.FinishActivation(activationId, status)
	if err != nil {
		if errors.Is(err, application.ActivationNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status": "ERROR",
				"error":  err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "SUCCESS",
	})
}
