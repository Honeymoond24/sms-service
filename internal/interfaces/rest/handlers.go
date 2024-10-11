package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
- GET_SERVICES - отдача списка сгененированных номеров
- GET_NUMBER - запрос номера по условиям (отдача случайного из сгененированных)
- FINISH_ACTIVATION - завершение работы с номером
- PUSH_SMS - отправка сгенерированной смс
*/

func FinishActivation(c echo.Context, body map[string]interface{}) error {
	c.Logger().Info("FinishActivation", body)
	// TODO
	return c.String(http.StatusOK, "Hello, World!")
}

func PushSms(c echo.Context, body map[string]interface{}) error {
	c.Logger().Info("PushSms", body)
	// TODO
	return c.String(http.StatusOK, "Hello, World!")
}
