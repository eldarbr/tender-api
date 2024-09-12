package repository

import (
	"avito-back-test/internal/db"
	"avito-back-test/internal/model"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoBid = errors.New("bid not found")
)

type BidRepository struct {
	db *sql.DB
}

func NewBidRepository() *BidRepository {
	db := db.DB
	return &BidRepository{
		db: db,
	}
}

func (r *BidRepository) InsertNewBid(b *model.Bid) error {
	bidQuery := `
INSERT INTO bid
	(tender_id, author_type, author_id)
VALUES ($1, $2, $3)
RETURNING 
	id,
	status,
	created_at
`
	bidInfoQuery := `
INSERT INTO bid_information
	(id, name, description)
VALUES ($1, $2, $3)
RETURNING
	version;
`
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(bidQuery, b.TenderID, b.AuthorType, b.AuthorID)
	err = row.Scan(&b.ID, &b.Status, &b.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	row = tx.QueryRow(bidInfoQuery, b.ID, b.Name, b.Description)
	err = row.Scan(&b.Version)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (r *BidRepository) GetUserBids(userID uuid.UUID, limit, offset int) ([]model.Bid, error) {
	query := `
SELECT
	b.id,
	bi.name,
	bi.description,
	b.tender_id,
	b.status,
	b.author_id,
	b.author_type,
	bi.version,
	b.created_at
FROM bid b
	JOIN bid_information bi
		ON bi.id = b.id
WHERE
	author_type = 'Organization'
	AND author_id IN (
		SELECT organization_id
		FROM organization_responsible
		WHERE user_id = $1
	)
	OR
	author_type = 'User'
	AND author_id = $1
ORDER BY name
LIMIT $2
OFFSET $3
`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []model.Bid
	for rows.Next() {
		var bid model.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description,
			&bid.TenderID, &bid.Status, &bid.AuthorID,
			&bid.AuthorType, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (r *BidRepository) GetBidsByTender(tenderID uuid.UUID, limit, offset int) ([]model.Bid, error) {
	query := `
SELECT
	b.id,
	bi.name,
	bi.description,
	b.tender_id,
	b.status,
	b.author_id,
	b.author_type,
	bi.version,
	b.created_at
FROM bid b
	JOIN bid_information bi
		ON bi.id = b.id
	JOIN (
		SELECT id, MAX(version) as mv
		FROM bid_information
		GROUP BY id
	) latest_bi
		ON latest_bi.id = bi.id AND bi.version = latest_bi.mv
WHERE
	b.tender_id = $1
ORDER BY name
LIMIT $2
OFFSET $3
`
	rows, err := r.db.Query(query, tenderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []model.Bid
	for rows.Next() {
		var bid model.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description,
			&bid.TenderID, &bid.Status, &bid.AuthorID,
			&bid.AuthorType, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (r *BidRepository) UpdateBidStatus(b *model.Bid) error {
	query := `
UPDATE bid
SET status = $1
WHERE
	id = $2
RETURNING
	id
`
	row := r.db.QueryRow(query, b.Status, b.ID)
	err := row.Scan(&b.ID)
	if err == sql.ErrNoRows {
		return ErrNoBid
	}
	if err != nil {
		return err
	}
	tmp, err := r.GetLastBidByID(b.ID)
	*b = *tmp
	return err
}

func (r *BidRepository) GetLastBidByID(bidID uuid.UUID) (*model.Bid, error) {
	query := `
SELECT
	b.id,
	bi.name,
	bi.description,
	b.tender_id,
	b.status,
	b.author_id,
	b.author_type,
	bi.version,
	b.created_at
FROM bid b
	JOIN bid_information bi
		ON bi.id = b.id
WHERE b.id = $1
ORDER BY version DESC
LIMIT 1
`
	var bid model.Bid

	row := r.db.QueryRow(query, bidID)
	err := row.Scan(&bid.ID, &bid.Name, &bid.Description,
		&bid.TenderID, &bid.Status, &bid.AuthorID,
		&bid.AuthorType, &bid.Version, &bid.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNoBid
	}
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *BidRepository) PatchBid(bidID uuid.UUID, patch *model.BidUpdate) (*model.Bid, error) {
	bidInfoQuery := `
INSERT INTO bid_information
	(id, name, description, version)
SELECT
	$1,
	$2,
	$3,
	$4,
	(SELECT MAX(version) + 1 FROM bid_information WHERE id = $1)
RETURNING
	version;
`
	b, err := r.GetLastBidByID(bidID)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		b.Name = *patch.Name
	}
	if patch.Description != nil {
		b.Description = *patch.Description
	}

	row := r.db.QueryRow(bidInfoQuery, b.ID, b.Name, b.Description)
	err = row.Scan(&b.Version)
	if err != nil {
		return nil, err
	}
	return r.GetLastBidByID(bidID)
}

func (r *BidRepository) RollbackBid(bidID uuid.UUID, version int) (*model.Bid, error) {
	bidInfoQuery := `
INSERT INTO bid_information
	(id, version, name, description)
SELECT
	id,
	(SELECT MAX(version) + 1 FROM bid_information WHERE id = $1),
	name,
	description,
	service_type
FROM bid_information
WHERE id = $1 AND version = $2
RETURNING version
`
	q, err := r.db.Query(bidInfoQuery, bidID, version)
	if err != nil {
		return nil, err
	}
	if !q.Next() {
		return nil, ErrNoBid
	}
	return r.GetLastBidByID(bidID)
}
