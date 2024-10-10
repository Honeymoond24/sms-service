package database

import (
	"database/sql"
	"github.com/Honeymoond24/sms-service/internal/domain"
)

type ServicesRepository struct {
	db *sql.DB
}

func NewServicesRepository(db *sql.DB) *ServicesRepository {
	return &ServicesRepository{db}
}

func (r *ServicesRepository) CreateService(service domain.Service) (int, error) {
	res, err := r.db.Exec("INSERT INTO services(name) VALUES(?)", service.Name)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}
