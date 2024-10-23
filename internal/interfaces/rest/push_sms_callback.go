package rest

import (
	"bytes"
	"encoding/json"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"log"
	"net/http"
)

func PushSms(url string, sms domain.SMS) {
	data := map[string]interface{}{
		"action":    "PUSH_SMS",
		"key":       "qwerty123", // TODO: get from env
		"smsID":     sms.ID,
		"phone":     sms.PhoneTo.Number,
		"phoneFrom": sms.PhoneFrom,
		"text":      sms.Text,
	}
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(data)
	if err != nil {
		log.Println("Error while encoding data: ", err)
		return
	}

	response, err := http.Post(url, "application/json", buff)
	if err != nil {
		log.Println("Error while sending request: ", err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Println("Error while closing response body: ", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		log.Println("Error while sending sms: ", response.Status)
		return
	}

	// unmarshal response
	var responseMap map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseMap)
	if err != nil {
		log.Println("Error while decoding PushSms response: ", err)
		return
	}
	log.Println("PushSms response: ", responseMap)
}
