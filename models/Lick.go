package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// Lick is used by pop to map your licks database table to your go code.
type Lick struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	TweetID      int64      `json:"tweet_id" db:"tweet_id"`
	TweetText    string     `json:"tweet_text" db:"tweet_text"`
	Bootlicker   Bootlicker `belongs_to:"bootlicker"`
	BootlickerID uuid.UUID  `json:"bootlicker_id" db:"bootlicker_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (l Lick) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// Licks is not required by pop and may be deleted
type Licks []Lick

// String is not required by pop and may be deleted
func (l Licks) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}
