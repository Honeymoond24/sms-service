package rest

import (
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"net/http"
)

func FinishActivation(c echo.Context, body map[string]interface{}, _ *application.SmsService) error {
	c.Logger().Info("FinishActivation", body)

	//activationIdFloat, ok := body["activationId"].(float64)
	//if !ok {
	//	return c.String(http.StatusBadRequest, "invalid smsID")
	//}
	//activationId := int(activationIdFloat)
	//
	//statusFloat, ok := body["smsID"].(float64)
	//if !ok {
	//	return c.String(http.StatusBadRequest, "invalid status")
	//}
	//status := int(statusFloat)

	return c.String(http.StatusOK, "Hello, World!")
}
