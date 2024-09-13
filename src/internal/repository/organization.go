package repository

import (
	"avito-back-test/internal/db"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrNoOrganization = errors.New("no organization with set uuid")
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository() *OrganizationRepository {
	db := db.DB
	return &OrganizationRepository{
		db: db,
	}
}

func (r *OrganizationRepository) GetOrganizationPresent(organizationID uuid.UUID) (bool, error) {
	query := `
SELECT 1
FROM organization
WHERE id = $1
`
	x, err := r.db.Query(query, organizationID)
	if err != nil {
		return false, err
	}
	if !x.Next() {
		return false, nil
	}
	return true, nil
}
