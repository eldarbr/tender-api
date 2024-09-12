package repository

import (
	"avito-back-test/internal/db"
	"avito-back-test/internal/model"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository() *EmployeeRepository {
	db := db.DB
	return &EmployeeRepository{
		db: db,
	}
}

var ErrNoEmployee = errors.New("no employees with set username")

func (r *EmployeeRepository) GetEmployeeByUsername(username string) (*model.Employee, error) {
	query := `
SELECT
	id,
	username,
	first_name,
	last_name,
	created_at,
	updated_at
FROM employee
WHERE username = $1`

	var employee model.Employee

	row := r.db.QueryRow(query, username)
	err := row.Scan(&employee.ID, &employee.Username, &employee.FirstName,
		&employee.LastName, &employee.CreatedAt, &employee.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrNoEmployee
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) GetEmployeeIDByUsername(username string) (*uuid.UUID, error) {
	employee, err := r.GetEmployeeByUsername(username)
	if err != nil {
		return nil, err
	}
	id := employee.ID
	return &id, nil
}

func (r *EmployeeRepository) GetEmployeePresent(employeeID uuid.UUID) (bool, error) {
	query := `
SELECT 1
FROM employee
WHERE id = $1
`
	x, err := r.db.Query(query, employeeID)
	if err != nil {
		return false, err
	}
	if !x.Next() {
		return false, nil
	}
	return true, nil
}

func (r *EmployeeRepository) GetEmployeeRespOrganization(employeeID uuid.UUID) (*uuid.UUID, error) {
	// Task says, one user is responsible for one organization at most
	query := `
SELECT organization_id
FROM organization_responsible
WHERE user_id = $1
`
	var organizationID uuid.UUID

	row := r.db.QueryRow(query, employeeID)
	err := row.Scan(&organizationID)

	if err == sql.ErrNoRows {
		return nil, ErrNoEmployee
	}
	if err != nil {
		return nil, err
	}
	return &organizationID, nil
}
