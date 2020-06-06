package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// PledgedDonation is used by pop to map your pledged_donations database table to your go code.
type PledgedDonation struct {
	ID uuid.UUID `json:"id" db:"id"`

	Amount       int       `json:"amount" db:"amount"`
	BootlickerID uuid.UUID `json:"bootlicker_id" db:"bootlicker_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p PledgedDonation) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PledgedDonations is not required by pop and may be deleted
type PledgedDonations []PledgedDonation

// String is not required by pop and may be deleted
func (p PledgedDonations) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}
