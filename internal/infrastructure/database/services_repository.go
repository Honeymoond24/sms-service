package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Honeymoond24/sms-service/internal/application"
	"github.com/Honeymoond24/sms-service/internal/domain"
	"maps"
	"strings"
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

func getServiceID(tx *sql.Tx, serviceCode string) (int, error) {
	row := tx.QueryRow(`SELECT id FROM services WHERE service_code = ?;`, serviceCode)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func buildLikeQueryArgs(column string, values []string) (string, []interface{}) {
	var placeholders []string
	var args []interface{}

	for _, value := range values {
		placeholders = append(placeholders, fmt.Sprintf("%s LIKE ?", column))
		args = append(args, value+"%")
	}
	query := strings.Join(placeholders, " OR ")

	return query, args
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
	country, service string,
	sum int,
	exceptionPhoneSet []string,
) (int, int, error) {
	fmt.Println("GetPhoneNumber")
	args := []interface{}{service, country}
	var queryArgs string
	if len(exceptionPhoneSet) > 0 {
		qa, queryArgsValues := buildLikeQueryArgs("p.number", exceptionPhoneSet)
		queryArgs = fmt.Sprintf(" AND (%s)", qa)
		args = append(args, queryArgsValues...)
	}
	query := `
		WITH sp AS (SELECT p.id
            FROM phone_numbers AS p
                     LEFT JOIN activations AS a ON p.id = a.phone_id
                     LEFT JOIN services AS s ON a.service_id = s.id
            WHERE s.service_code = ?)
		SELECT p.id, p.number
		FROM phone_numbers AS p
				 JOIN countries AS c ON p.country_id = c.id
				 LEFT JOIN activations AS a ON p.id = a.phone_id
		WHERE c.name = ? AND p.id NOT IN (SELECT id FROM sp)
	` + queryArgs + ` LIMIT 1;`

	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println("Error while starting transaction:", err)
		return 0, 0, err
	}
	row := tx.QueryRow(query, args...)

	var phoneId, number int
	err = row.Scan(&phoneId, &number)
	if err != nil {
		fmt.Println("Error while scanning rows:", err)
		return 0, 0, err
	}

	serviceIdCh := make(chan int)
	go func(service string, tx *sql.Tx, serviceIdCh chan int) {
		id, err := getServiceID(tx, service)
		if err != nil {
			fmt.Println("Error while getting service ID:", err)
			return
		}
		serviceIdCh <- id
	}(service, tx, serviceIdCh)

	var activationID int
	row = tx.QueryRow(`
		INSERT INTO activations (sum_price, status, phone_id, service_id) 
		VALUES (?, ?, ?, ?) RETURNING id;`, sum, 1, phoneId, <-serviceIdCh)
	err = row.Scan(&activationID)
	if err != nil {
		fmt.Println("Error while inserting activation:", err)
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error while committing transaction:", err)
		return 0, 0, err
	}

	return number, activationID, nil
}

func (r *SMSServiceRepository) StoreSms(sms domain.SMS) error {
	_, err := r.db.Exec(`
		INSERT INTO sms (sms_id, phone_id, phone_from, text)
		VALUES (?, ?, ?, ?);
	`, sms.ID, sms.PhoneTo.ID, sms.PhoneFrom, sms.Text)
	if err != nil {
		return err
	}

	return nil
}

func (r *SMSServiceRepository) GetPhoneNumberByPhone(phone int) (domain.PhoneNumber, error) {
	var number domain.PhoneNumber
	err := r.db.QueryRow(`
		SELECT id
		FROM phone_numbers
		WHERE number = ?;
	`, phone).Scan(&number.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PhoneNumber{}, nil
		}
		return domain.PhoneNumber{}, err
	}
	number.Number = phone

	return number, nil
}

func (r *SMSServiceRepository) FinishActivation(activationId, status int) error {
	result, err := r.db.Exec(`
		UPDATE activations
		SET status = ?
		WHERE id = ?;
	`, status, activationId)
	if err != nil {
		fmt.Println("Error while updating activation status:", err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		fmt.Println("Error while getting rows affected:", err)
		return err
	} else {
		if rowsAffected == 0 {
			return application.ActivationNotFound
		}
	}
	return nil
}
