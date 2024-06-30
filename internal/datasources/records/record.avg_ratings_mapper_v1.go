package records

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
)

func (c *AvgRatings) ToV1Domain() V1Domains.AvgRatingDomain {
	return V1Domains.AvgRatingDomain{
		ID:           c.ID,
		ProductID:    c.ProductID,
		Overall:      c.Overall,
		Quality:      c.Quality,
		Value:        c.Value,
		TotalReviews: c.TotalReviews,
		Star5:        c.Star5,
		Star4:        c.Star4,
		Star3:        c.Star3,
		Star2:        c.Star2,
		Star1:        c.Star1,
		TimeBaseDomain: V1Domains.TimeBaseDomain{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		},
	}
}

func FromAvgRatingV1Domain(c *V1Domains.AvgRatingDomain) AvgRatings {
	return AvgRatings{
		ID:           c.ID,
		ProductID:    c.ProductID,
		Overall:      c.Overall,
		Quality:      c.Quality,
		Value:        c.Value,
		TotalReviews: c.TotalReviews,
		Star5:        c.Star5,
		Star4:        c.Star4,
		Star3:        c.Star3,
		Star2:        c.Star2,
		Star1:        c.Star1,
		TimeBase: TimeBase{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		},
	}
}

func ToAvgRatingV1Domain(c *[]AvgRatings) []V1Domains.AvgRatingDomain {
	var result []V1Domains.AvgRatingDomain

	for _, val := range *c {
		result = append(result, val.ToV1Domain())
	}

	return result
}
