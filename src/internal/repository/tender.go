package repository

import (
	"avito-back-test/internal/db"
	"avito-back-test/internal/model"
	"database/sql"
)

type TenderRepository struct {
	db *sql.DB
}

func NewTenderRepository() *TenderRepository {
	db := db.DB
	return &TenderRepository{
		db: db,
	}
}

func (r *TenderRepository) GetAllTenders() ([]model.Tender, error) {

	query := `SELECT id, name, description, service_type, status, organization_id, version, created_at FROM tender`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []model.Tender
	for rows.Next() {
		var tender model.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description,
			&tender.ServiceType, &tender.Status, &tender.OrganizationID,
			&tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}
	return tenders, nil
}
