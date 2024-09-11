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

func (r *TenderRepository) GetAllPublicTenders(limit, offset int) ([]model.Tender, error) {
	query := `
SELECT
	t.id,
	ti.name,
	ti.description,
	ti.service_type,
	t.status,
	t.organization_id,
	ti.version,
	t.created_at
FROM tender t
	JOIN tender_information ti
		ON ti.id = t.id
	JOIN (
		SELECT id, MAX(version) as mv
		FROM tender_information
		GROUP BY id
	) latest_ti
		ON latest_ti.id = ti.id AND ti.version = latest_ti.mv
WHERE status = 'Published'
ORDER BY name
LIMIT $1
OFFSET $2
`
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
			&tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}
	return tenders, nil
}

func (r *TenderRepository) GetPublicTendersOfService(serviceType string, limit, offset int) ([]model.Tender, error) {
	query := `
SELECT
	t.id,
	ti.name,
	ti.description,
	ti.service_type,
	t.status,
	t.organization_id,
	ti.version,
	t.created_at
FROM tender t
	JOIN tender_information ti
		ON ti.id = t.id
	JOIN (
		SELECT id, MAX(version) as mv
		FROM tender_information
		GROUP BY id
	) latest_ti
		ON latest_ti.id = ti.id AND ti.version = latest_ti.mv
WHERE
	status = 'Published'
	AND service_type = $1
ORDER BY name
LIMIT $2
OFFSET $3
`
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
			&tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}
	return tenders, nil
}

func (r *TenderRepository) InsertNewTender(t *model.Tender) error {
	tenderQuery := `
INSERT INTO tender
	(organization_id)
VALUES ($1)
RETURNING 
	id,
	status,
	created_at
`
	tenderInfoQuery := `
INSERT INTO tender_information
	(id, name, description, service_type)
VALUES
	($1, $2, $3, $4)
RETURNING
	version;
`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(tenderQuery, t.OrganizationID)
	err = row.Scan(&t.ID, &t.Status, &t.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	row = tx.QueryRow(tenderInfoQuery, t.ID, t.Name, t.Description, t.ServiceType)
	err = row.Scan(&t.Version)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (r *TenderRepository) GetUserTenders(userID uuid.UUID, limit, offset int) ([]model.Tender, error) {
	query := `
SELECT
	t.id,
	ti.name,
	ti.description,
	ti.service_type,
	t.status,
	t.organization_id,
	ti.version,
	t.created_at
FROM tender t
	JOIN tender_information ti
		ON ti.id = t.id
WHERE organization_id IN (
	SELECT organization_id
	FROM organization_responsible
	WHERE user_id = $1
)
ORDER BY name
LIMIT $2
OFFSET $3
`
	rows, err := r.db.Query(query, userID, limit, offset)
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
WHERE
	id = $2
RETURNING
	id
`
	row := r.db.QueryRow(query, t.Status, t.ID)
	err := row.Scan(&t.ID)
	if err == sql.ErrNoRows {
		return ErrNoTender
	}
	if err != nil {
		return err
	}
	tmp, err := r.GetLastTenderByID(t.ID)
	*t = *tmp
	return err
}

func (r *TenderRepository) GetLastTenderByID(tenderID uuid.UUID) (*model.Tender, error) {
	query := `
SELECT
	t.id,
	ti.name,
	ti.description,
	ti.service_type,
	t.status,
	t.organization_id,
	ti.version,
	t.created_at
FROM tender t
	JOIN tender_information ti
		ON ti.id = t.id
WHERE t.id = $1
ORDER BY version DESC
LIMIT 1
`

	var t model.Tender

	row := r.db.QueryRow(query, tenderID)
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNoTender
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TenderRepository) PatchTender(tenderID uuid.UUID, patch *model.TenderUpdate) (*model.Tender, error) {
	tenderInfoQuery := `
INSERT INTO tender_information
	(id, name, description, service_type, version)
SELECT
	$1,
	$2,
	$3,
	$4,
	(SELECT MAX(version) + 1 FROM tender_information WHERE id = $1)
RETURNING
	version;
`
	t, err := r.GetLastTenderByID(tenderID)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		t.Name = *patch.Name
	}
	if patch.Description != nil {
		t.Description = *patch.Description
	}
	if patch.ServiceType != nil {
		t.ServiceType = *patch.ServiceType
	}

	row := r.db.QueryRow(tenderInfoQuery, t.ID, t.Name, t.Description, t.ServiceType)
	err = row.Scan(&t.Version)
	if err != nil {
		return nil, err
	}
	return r.GetLastTenderByID(tenderID)
}

func (r *TenderRepository) RollbackTender(tenderID uuid.UUID, version int) (*model.Tender, error) {
	tenderInfoQuery := `
INSERT INTO tender_information
	(id, version, name, description, service_type)
SELECT
	id,
	(SELECT MAX(version) + 1 FROM tender_information WHERE id = $1),
	name,
	description,
	service_type
FROM tender_information
WHERE id = $1 AND version = $2
RETURNING version
`
	q, err := r.db.Query(tenderInfoQuery, tenderID, version)
	if err != nil {
		return nil, err
	}
	if !q.Next() {
		return nil, ErrNoTender
	}
	return r.GetLastTenderByID(tenderID)
}
