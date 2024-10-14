package rest

import (
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type country struct {
	Country     string                    `json:"country"`
	OperatorMap map[string]map[string]int `json:"operatorMap"`
}
type GetServicesResponse struct {
	CountryList []country `json:"countryList"`
	Status      string    `json:"status"`
}

func GetServices(c echo.Context, body map[string]interface{}, service *application.SmsService) error {
	c.Logger().Info("GetServices", body)

	countries, err := service.GetServices()
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	var countryList []country
	for countryName, services := range countries {
		c := country{
			Country:     countryName,
			OperatorMap: make(map[string]map[string]int),
		}
		c.OperatorMap["any"] = make(map[string]int)
		for serviceCode, diff := range services {
			c.OperatorMap["any"][serviceCode] = diff
		}

		countryList = append(countryList, c)
	}

	response := GetServicesResponse{
		CountryList: countryList,
		Status:      "SUCCESS",
	}

	return c.JSON(http.StatusOK, response)
}
