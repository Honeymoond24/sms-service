package rest

import (
	"fmt"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetNumber(c echo.Context, body map[string]interface{}, service *application.SmsService) error {
	c.Logger().Info("GetNumber", body)

	countryName, ok := body["country"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "Invalid country")
	}

	serviceName, ok := body["service"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "Invalid service")
	}

	sum, ok := body["sum"].(float64)
	if !ok {
		return c.String(http.StatusBadRequest, "Invalid sum")
	}

	var phonePrefixes []string
	if exceptions, ok := body["exceptionPhoneSet"].([]interface{}); ok {
		phonePrefixes = make([]string, len(exceptions))
		for i, v := range exceptions {
			if str, ok := v.(string); ok {
				phonePrefixes[i] = str
			} else {
				return c.String(http.StatusBadRequest, "Invalid phone prefix")
			}
		}
	}

	number, activationID, err := service.GetNumber(
		countryName,
		serviceName,
		int(sum),
		phonePrefixes,
	)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	if number == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "NO_NUMBERS",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":        "SUCCESS",
		"number":        number,
		"activation_id": activationID,
	})
}
