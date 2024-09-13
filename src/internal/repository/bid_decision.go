package repository

import (
	"avito-back-test/internal/db"
	"database/sql"
	"github.com/google/uuid"
)

type BidDecisionRepository struct {
	db *sql.DB
}

func NewBidDecisionRepository() *BidDecisionRepository {
	db := db.DB
	return &BidDecisionRepository{
		db: db,
	}
}

func (r *BidDecisionRepository) TxInsertUpdateDecision(tx *sql.Tx, bidID, userID uuid.UUID, decision string) error {
	insertQuery := `
INSERT INTO bid_decision
	(bid_id, responsible_id, decision)
VALUES
	($1, $2, $3)
`
	updateQuery := `
UPDATE bid_decision
SET decision = $3
WHERE bid_id = $1 AND responsible_id = $2
`
	res, err := tx.Exec(updateQuery, bidID, userID, decision)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		_, err = tx.Exec(insertQuery, bidID, userID, decision)
	}
	return err
}

func (r *BidDecisionRepository) TxCountDecisions(tx *sql.Tx, bidID uuid.UUID) (int, int, error) {
	query := `
SELECT
	approve.approve_count, reject.reject_count
FROM
(
	SELECT COUNT(1) as approve_count
	FROM bid_decision
	WHERE bid_id = $1 AND decision = 'Approved'
) as approve,
(
	SELECT COUNT(1) as reject_count
	FROM bid_decision
	WHERE bid_id = $1 AND decision = 'Rejected'
) as reject
`
	var (
		approveCount int
		rejectCount  int
	)

	row := tx.QueryRow(query, bidID)
	err := row.Scan(&approveCount, &rejectCount)
	return approveCount, rejectCount, err
}

func (r *BidDecisionRepository) WithTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)

	return err
}
