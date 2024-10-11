package rest

import (
	"encoding/json"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Server struct {
	srv          *echo.Echo
	servicesRepo application.ServicesRepository
}

func NewServer(repo application.ServicesRepository) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	s := &Server{srv: e, servicesRepo: repo}
	e.POST("/", s.HandlerBase)
	return s
}

func (s *Server) Run() {
	s.srv.Logger.Fatal(s.srv.Start(":8080"))
}

func (s *Server) HandlerBase(c echo.Context) error {
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Bad request")
	}
	if _, ok := jsonMap["action"]; !ok {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	if _, ok := jsonMap["key"]; !ok {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	switch jsonMap["action"] {
	case "GET_SERVICES":
		return GetServices(c, jsonMap, s.servicesRepo)
	case "GET_NUMBER":
		return GetNumber(c, jsonMap, s.servicesRepo)
	case "FINISH_ACTIVATION":
		return FinishActivation(c, jsonMap)
	case "PUSH_SMS":
		return PushSms(c, jsonMap)
	default:
		return c.String(http.StatusBadRequest, "Bad request")
	}
}
