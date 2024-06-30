package records

import (
	"time"
)

type TimeBase struct {
	CreatedAt time.Time `db:"created_at"`
	// It is also a pointer to a time.Time value,
	// allowing it to be nil if the record is not deleted.
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type DetailBase struct {
	Detail map[string][]string `json:"detail"`
}

type SpecBase struct {
	Spec map[string]string `json:"specification"`
}
