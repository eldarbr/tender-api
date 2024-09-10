package model

import (
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	ID          uuid.UUID
	Name        string
	Description string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
