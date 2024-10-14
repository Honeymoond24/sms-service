package rest

import (
	"encoding/json"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Server struct {
	srv     *echo.Echo
	service *application.SmsService
}

func NewServer(service *application.SmsService) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	s := &Server{srv: e, service: service}
	e.POST("/", s.HandlerBase)
	return s
}

func (s *Server) Run() {
	s.srv.Logger.Fatal(s.srv.Start("127.0.0.1:8000"))
}

func (s *Server) HandlerBase(c echo.Context) error {
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusUnprocessableEntity, "unprocessable entity")
	}
	if _, ok := jsonMap["action"]; !ok {
		return c.String(http.StatusBadRequest, "action field required")
	}
	if _, ok := jsonMap["key"]; !ok {
		return c.String(http.StatusBadRequest, "key field required")
	}

	switch jsonMap["action"] {
	case "GET_SERVICES":
		return GetServices(c, jsonMap, s.service)
	case "GET_NUMBER":
		return GetNumber(c, jsonMap, s.service)
	case "PUSH_SMS":
		return PushSms(c, jsonMap, s.service)
	case "FINISH_ACTIVATION":
		return FinishActivation(c, jsonMap, s.service)
	default:
		return c.String(http.StatusBadRequest, "unknown action type")
	}
}
