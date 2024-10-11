package database

import (
	"database/sql"
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

func (r *SMSServiceRepository) GetPhoneNumber(
	country, _ string, // country, service
	_ int, // sum
	exceptionPhoneSet []string,
) (string, int, error) {
	// TODO: Implement filtering by exceptionPhoneSet
	args := []interface{}{country}
	if len(exceptionPhoneSet) > 0 {
		args = append(args, fmt.Sprintf("%s%%", exceptionPhoneSet[0]))
	}
	rows, err := r.db.Query(`
		SELECT pn.number, a.id
		FROM phone_numbers AS pn
		JOIN countries AS c ON pn.country_id = c.id
		JOIN activations AS a ON pn.activation_id = a.id
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
