package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"maps"
)

type SMSServiceRepository struct {
	db *sql.DB
}

func NewServicesRepository(db *sql.DB) *SMSServiceRepository {
	return &SMSServiceRepository{db}
}

func (r *SMSServiceRepository) GetServiceCodes(ch chan map[string]int) error {
	rows, err := r.db.Query(`SELECT service_code FROM services;`)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error while closing rows:", err)
		}
	}(rows)

	serviceCodes := make(map[string]int)
	for rows.Next() {
		var code string
		err := rows.Scan(&code)
		if err != nil {
			return err
		}
		serviceCodes[code] = 0
	}

	ch <- serviceCodes
	return nil
}

func (r *SMSServiceRepository) GetServices() (map[string]map[string]int, error) {
	serviceCodesCh := make(chan map[string]int)
	go func(ch chan map[string]int) {
		err := r.GetServiceCodes(ch)
		if err != nil {
			fmt.Println("Error while getting service codes:", err)
		}
	}(serviceCodesCh)
	rows, err := r.db.Query(`
		WITH phone_count AS (SELECT c.name AS country, COUNT(p.id) AS count
                     FROM phone_numbers p
                              JOIN countries c on c.id = p.country_id
                     GROUP BY c.name)
		SELECT
			countries.name AS country,
			ifnull(services.service_code, '') AS service,
			pc.count - COUNT(activations.id) AS diff,
			pc.count AS total
		FROM phone_numbers
				 LEFT JOIN activations ON phone_numbers.id = activations.phone_id
				 LEFT JOIN countries ON phone_numbers.country_id = countries.id
				 LEFT JOIN services ON activations.service_id = services.id
				JOIN phone_count pc ON pc.country = countries.name
		GROUP BY countries.name, services.service_code;
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error while closing rows:", err)
		}
	}(rows)

	serviceCodes := <-serviceCodesCh
	fmt.Println(serviceCodes)

	countries := make(map[string]map[string]int)
	for rows.Next() {
		var country, serviceCode string
		var diff, total int
		serviceCodesPerCountry := make(map[string]int)
		maps.Copy(serviceCodesPerCountry, serviceCodes)

		err := rows.Scan(&country, &serviceCode, &diff, &total)

		for k, _ := range serviceCodesPerCountry {
			serviceCodesPerCountry[k] = total
		}
		fmt.Println("serviceCodesPerCountry", serviceCodesPerCountry)
		if _, ok := countries[country]; ok {
			countries[country][serviceCode] = diff
		} else {
			countries[country] = serviceCodesPerCountry
		}
		if err != nil {
			return nil, err
		}
	}

	return countries, nil
}

// GetPhoneNumber
// Input example:
//
//	"country": "russia",
//	"operator": "any",
//	"service": "tg",
//	"sum": 20.00,
/*
Actions:
- Get number from DB by filters: where phone number and service not in activations and provided country match
- Save sum, service, phone number into Activations
- return number and activation ID
*/
func (r *SMSServiceRepository) GetPhoneNumber(
	country, service string, // country, service
	sum int, // sum
	exceptionPhoneSet []string,
) (string, int, error) {
	fmt.Println("GetPhoneNumber", country, service, sum, exceptionPhoneSet)
	// TODO: Implement filtering by exceptionPhoneSet
	args := []interface{}{country}
	if len(exceptionPhoneSet) > 0 {
		args = append(args, fmt.Sprintf("%s%%", exceptionPhoneSet[0]))
	}
	rows, err := r.db.Query(`
		SELECT pn.number
		FROM phone_numbers AS pn
		JOIN countries AS c ON pn.country_id = c.id
		WHERE c.name = ? AND pn.number NOT LIKE ?;
	`, args...)
	if err != nil {
		return "", 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error while closing rows:", err)
		}
	}(rows)

	var number string
	var activationID int
	for rows.Next() {
		err := rows.Scan(&number, &activationID)
		if err != nil {
			return "", 0, err
		}
	}

	return number, activationID, nil
}

func (r *SMSServiceRepository) StoreSms(sms domain.SMS) error {
	_, err := r.db.Exec(`
		INSERT INTO sms (sms_id, phone, phone_from, text)
		VALUES (?, ?, ?, ?);
	`, sms.ID, sms.Phone, sms.PhoneFrom, sms.Text)
	if err != nil {
		return err
	}

	return nil
}

func (r *SMSServiceRepository) GetPhoneNumberByPhone(phone int) (domain.PhoneNumber, error) {
	var number domain.PhoneNumber
	err := r.db.QueryRow(`
		SELECT id, number
		FROM phone_numbers
		WHERE number = ?;
	`, phone).Scan(&number.ID, &number.Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PhoneNumber{}, nil
		}
		return domain.PhoneNumber{}, err
	}

	return number, nil
}
