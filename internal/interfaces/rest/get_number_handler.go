package rest

import (
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetNumber(c echo.Context, body map[string]interface{}, repo application.ServicesRepository) error {
	c.Logger().Info("GetNumber", body)

	country, ok := body["country"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "Invalid country")
	}

	service, ok := body["service"].(string)
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

	number, activationID, err := repo.GetPhoneNumber(
		country,
		service,
		int(sum),
		phonePrefixes,
	)
	if err != nil {
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
