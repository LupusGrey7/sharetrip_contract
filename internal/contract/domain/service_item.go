package domain

import "time"

// OfferingEntity — dictionary entity used by usecase and storage.
type OfferingEntity struct {
	ServiceCode string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
