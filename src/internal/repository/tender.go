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
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at
FROM tender
WHERE status = 'Published'
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
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at
FROM tender
WHERE
	service_type = $1
	AND status = 'Published'
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
			&tender.Version, &tender.CreatedAt)
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
	(name, description, service_type, organization_id)
VALUES
	($1, $2, $3, $4)
RETURNING
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at`

	row := r.db.QueryRow(query, t.Name, t.Description, t.ServiceType, t.OrganizationID)
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt)
	return err
}

func (r *TenderRepository) GetUserTenders(userID uuid.UUID, limit, offset int) ([]model.Tender, error) {
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
WHERE organization_id IN (
	SELECT organization_id
	FROM organization_responsible
	WHERE user_id = $1
)
ORDER BY name
LIMIT $2
OFFSET $3`
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
	AND version = $3
RETURNING
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at`

	row := r.db.QueryRow(query, t.Status, t.ID, t.Version)
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

func (r *TenderRepository) GetLastTenderByID(tenderID uuid.UUID) (*model.Tender, error) {
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
WHERE id = $1
ORDER BY version DESC
LIMIT 1`

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
	closeQuery := `
UPDATE tender
SET status = 'Closed'
WHERE
	id = $1
	AND version + 1 = (
		SELECT MAX(version)
		FROM tender
		WHERE id = $1
	)
`
	insertQuery := `
WITH last_version AS (
	SELECT *
	FROM tender
	WHERE id = $1
	ORDER BY version DESC
	LIMIT 1
)
INSERT INTO tender
	(id, name, description, service_type, status, organization_id, version)
SELECT
	id,
	COALESCE($2, (SELECT name FROM last_version)),
	COALESCE($3, (SELECT description FROM last_version)),
	COALESCE($4, (SELECT service_type FROM last_version)),
	status,
	organization_id,
	version + 1
FROM last_version
RETURNING
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at;
`

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var t model.Tender

	row := tx.QueryRow(insertQuery, tenderID, patch.Name, patch.Description, patch.ServiceType)
	err = row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(closeQuery, tenderID)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &t, err
}

func (r *TenderRepository) RollbackTender(tenderID uuid.UUID, version int) (*model.Tender, error) {
	insertQuery := `
WITH max_version AS (
	SELECT MAX(version) AS value
	FROM tender
	WHERE id = $1
)
INSERT INTO tender
	(id, name, description, service_type, status, organization_id, version)
SELECT
	id,
	name,
	description,
	service_type,
	(
		SELECT status
		FROM tender
		WHERE
			id = $1
			AND version = max_version.value
	),
	organization_id,
	max_version.value + 1
FROM tender, max_version
WHERE
	id = $1
	AND version = $2
RETURNING
	id,
	name,
	description,
	service_type,
	status,
	organization_id,
	version,
	created_at;
`
	closeQuery := `
UPDATE tender
SET status = 'Closed'
WHERE
	id = $1
	AND version + 1 = (
		SELECT MAX(version)
		FROM tender
		WHERE id = $1
	)
`

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var t model.Tender

	row := tx.QueryRow(insertQuery, tenderID, version)
	err = row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status,
		&t.OrganizationID, &t.Version, &t.CreatedAt)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			err = ErrNoTender
		}
		return nil, err
	}

	_, err = tx.Exec(closeQuery, tenderID)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return &t, err
}
