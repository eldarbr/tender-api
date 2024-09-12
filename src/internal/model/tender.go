package model

import (
	"github.com/google/uuid"
	"time"
)

type Tender struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	ServiceType    string    `json:"serviceType"`
	Version        int       `json:"version"`
	OrganizationID uuid.UUID `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
}

type TenderUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ServiceType *string `json:"serviceType,omitempty"`
}
