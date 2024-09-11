package repository

import (
	"avito-back-test/internal/db"
	"avito-back-test/internal/model"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type TenderRepository struct {
	db *sql.DB
}

var ErrNoTender = errors.New("tender with set id not found")

func NewTenderRepository() *TenderRepository {
	db := db.DB
	return &TenderRepository{
		db: db,
	}
}

func (r *TenderRepository) GetAllTenders(limit, offset int) ([]model.Tender, error) {
	query := `
SELECT
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at,
	creator_username
FROM tender
ORDER BY name
LIMIT $1
OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []model.Tender
	for rows.Next() {
		var tender model.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description,
			&tender.ServiceType, &tender.Status, &tender.OrganizationID,
			&tender.Version, &tender.CreatedAt, &tender.CreatorUsername)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}
	return tenders, nil
}

func (r *TenderRepository) GetTendersOfService(serviceType string, limit, offset int) ([]model.Tender, error) {
	query := `
SELECT
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at,
	creator_username
FROM tender
WHERE service_type = $1
ORDER BY name
LIMIT $2
OFFSET $3`
	rows, err := r.db.Query(query, serviceType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []model.Tender
	for rows.Next() {
		var tender model.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description,
			&tender.ServiceType, &tender.Status, &tender.OrganizationID,
			&tender.Version, &tender.CreatedAt, &tender.CreatorUsername)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}
	return tenders, nil
}

func (r *TenderRepository) InsertNewTender(t *model.Tender) error {
	query := `
INSERT INTO tender
	(name, description, service_type, organization_id, creator_username)
VALUES
	($1, $2, $3, $4, $5)
RETURNING
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at`

	row := r.db.QueryRow(query, t.Name, t.Description, t.ServiceType, t.OrganizationID, t.CreatorUsername)
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt)
	return err
}

func (r *TenderRepository) GetUserTenders(username string, limit, offset int) ([]model.Tender, error) {
	query := `
SELECT
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at
FROM tender
WHERE creator_username = $1
ORDER BY name
LIMIT $2
OFFSET $3`
	rows, err := r.db.Query(query, username, limit, offset)
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

func (r *TenderRepository) UpdateTenderStatus(t *model.Tender) error {
	query := `
UPDATE tender
SET status = $1
WHERE id = $2
RETURNING
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at`

	row := r.db.QueryRow(query, t.Status, t.ID)
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return ErrNoTender
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *TenderRepository) GetTenderByID(tenderID uuid.UUID) (*model.Tender, error) {
	query := `
SELECT
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at,
	creator_username
FROM tender
WHERE id = $1`

	var t model.Tender

	row := r.db.QueryRow(query, tenderID)
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt, &t.CreatorUsername)
	if err == sql.ErrNoRows {
		return nil, ErrNoTender
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
