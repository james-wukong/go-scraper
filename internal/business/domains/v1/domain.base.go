package v1

import "time"

type TimeBaseDomain struct {
	CreatedAt time.Time
	// It is also a pointer to a time.Time value,
	// allowing it to be nil if the record is not deleted.
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type DetailBaseDomain struct {
	Detail map[string][]string
}

type SpecBaseDomain struct {
	Spec map[string]string
}
