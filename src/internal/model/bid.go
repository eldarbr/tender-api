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
	CreatedAt   time.Time `json:"createdAt"`
}

type BidUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type BidReview struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type BidStatus = string

const (
	BidCreated   BidStatus = "Created"
	BidPublished BidStatus = "Published"
	BidCanceled  BidStatus = "Canceled"
)

type AuthorType = string

const (
	AuthorTypeUser         AuthorType = "User"
	AuthorTypeOrganization AuthorType = "Organization"
)

type BidDecisionType = string

const (
	BidDecisionApproved BidDecisionType = "Approved"
	BidDecisionRejected BidDecisionType = "Rejected"
)
