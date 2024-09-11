package repository

import (
	"avito-back-test/internal/db"
	"database/sql"

	"github.com/google/uuid"
)

type OrganizationResponsibleRepository struct {
	db *sql.DB
}

func NewOrganizationResponsibleRepository() *OrganizationResponsibleRepository {
	db := db.DB
	return &OrganizationResponsibleRepository{
		db: db,
	}
}

func (r *OrganizationResponsibleRepository) GetIfEmployeeIsResponsible(employeeID, organizationID *uuid.UUID) (bool, error) {
	query := `
SELECT
	id
FROM organization_responsible
WHERE
	user_id = $1
	AND organization_id = $2`

	row := r.db.QueryRow(query, employeeID, organizationID)
	var id uuid.UUID
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
