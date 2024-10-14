package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Honeymoond24/sms-service/internal/domain"
)

type SMSServiceRepository struct {
	db *sql.DB
}

func NewServicesRepository(db *sql.DB) *SMSServiceRepository {
	return &SMSServiceRepository{db}
}

func (r *SMSServiceRepository) GetServices() (map[string][]domain.Service, error) {
	rows, err := r.db.Query(`
		SELECT c.name, s.name, s.amount
		FROM countries AS c
		JOIN services AS s ON c.id = s.country_id;
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

	countries := make(map[string][]domain.Service)
	for rows.Next() {
		var service domain.Service
		var country string
		err := rows.Scan(&country, &service.Name, &service.Amount)
		if services, ok := countries[country]; ok {
			countries[country] = append(services, service)
		} else {
			countries[country] = []domain.Service{service}
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
