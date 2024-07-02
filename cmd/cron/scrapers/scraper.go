package scrapers

import (
	"time"
)

type Prices struct {
	InWarehouse float32
	EcoFee      float32
	InstSave    float32
	Price       float32
	*Valid
}

type Valid struct {
	StartAt  time.Time
	EndedAt  time.Time
	Duration string
}

type FBPrices struct {
	CategoryUrl string
	Prices
}
