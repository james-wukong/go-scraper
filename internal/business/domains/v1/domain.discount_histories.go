package v1

import (
	"context"
	"time"
)

type DiscHistoryDomain struct {
	ID          uint
	ProductID   string
	Price       float32
	SaveAmount  float32
	SavePercent float32
	Duration    string
	StartedAt   time.Time
	EndedAt     time.Time
	TimeBaseDomain
}

type DiscHistoryUCInterface interface {
	// used by handler
	Store(inDom *DiscHistoryDomain) (outDom DiscHistoryDomain, statusCode int, err error)
	// used by handler
	GetByName(ctx context.Context, name string) (outDom DiscHistoryDomain, statusCode int, err error)
}

type DiscHistoryRepoInterface interface {
	// used by scrapper: get product id by it's name and sku
	GetHistoryByPID(pid string) (outDom DiscHistoryDomain, err error)
	// used by scrapper
	SaveDiscHistory(inDom *DiscHistoryDomain) (lastInsertId uint, err error)
	// used by handler
	// GetByName(name string) (outDom CategoryDomain, err error)
}
