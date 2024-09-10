package repository

import (
	// "avito-back-test/internal/model"
	"avito-back-test/internal/db"
	"database/sql"
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
