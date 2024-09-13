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

func (r *OrganizationResponsibleRepository) TxGetResponsibleCountByEmployee(tx *sql.Tx, employeeID uuid.UUID) (int, error) {
	query := `
SELECT COUNT(1)
FROM organization_responsible ores
	JOIN (
		SELECT organization_id
		FROM organization_responsible
		WHERE user_id = $1
	) org
		ON org.organization_id = ores.organization_id
`
	var count int
	row := tx.QueryRow(query, employeeID)
	err := row.Scan(&count)
	if err != nil {
		return 0, nil
	}
	if count == 0 {
		return 0, ErrNoOrganization
	}
	return count, nil
}
