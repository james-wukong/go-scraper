package v1

import (
	"context"
)

type AvgRatingDomain struct {
	ID           uint
	ProductID    string
	Overall      float32
	Quality      float32
	Value        float32
	TotalReviews int
	Star5        int
	Star4        int
	Star3        int
	Star2        int
	Star1        int
	TimeBaseDomain
}

type AvgRatingUCInterface interface {
	// used by scrapper
	Store(inDom *AvgRatingDomain) (outDom AvgRatingDomain, statusCode int, err error)
	// used by handler
	GetByName(ctx context.Context, name string) (outDom AvgRatingDomain, statusCode int, err error)
}

type AvgRatingRepoInterface interface {
	// used by scrapper: get categoryId by it's name and platform
	GetAvgRatingByPID(productID string) (outDom AvgRatingDomain, err error)
	// used by scrapper
	UpsertAvgRating(inDom *AvgRatingDomain) (lastInsertId uint, err error)
	// used by handler
	// GetByName(name string) (outDom CategoryDomain, err error)
}
