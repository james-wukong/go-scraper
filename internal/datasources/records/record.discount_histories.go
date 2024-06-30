package records

import "time"

type DiscountHistories struct {
	ID          uint      `db:"id"`
	ProductID   string    `db:"product_id"`
	Price       float32   `db:"price"`
	SaveAmount  float32   `db:"save_amount"`
	SavePercent float32   `db:"save_percent"`
	Duration    string    `db:"duration"`
	StartedAt   time.Time `db:"started_at"`
	EndedAt     time.Time `db:"ended_at"`
	TimeBase
}
