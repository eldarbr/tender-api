package model

import (
	"github.com/google/uuid"
	"time"
)

type Bid struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TenderID    uuid.UUID `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorID    uuid.UUID `json:"authorId"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
}

type BidUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
