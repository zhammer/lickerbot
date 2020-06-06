package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// Bootlicker is used by pop to map your bootlickers database table to your go code.
type Bootlicker struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	TwitterUserID    int64            `json:"twitter_user_id" db:"twitter_user_id"`
	TwitterHandle    string           `json:"twitter_handle" db:"twitter_handle"`
	Licks            Licks            `has_many:"licks" order_by:"id"`
	PledgedDonations PledgedDonations `has_many:"pledged_donations" order_by:"id"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (b Bootlicker) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Bootlickers is not required by pop and may be deleted
type Bootlickers []Bootlicker

// String is not required by pop and may be deleted
func (b Bootlickers) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}
