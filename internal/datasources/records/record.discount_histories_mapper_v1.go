package records

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
)

func (c *DiscountHistories) ToV1Domain() V1Domains.DiscHistoryDomain {
	return V1Domains.DiscHistoryDomain{
		ID:          c.ID,
		ProductID:   c.ProductID,
		Price:       c.Price,
		SaveAmount:  c.SaveAmount,
		SavePercent: c.SavePercent,
		Duration:    c.Duration,
		StartedAt:   c.StartedAt,
		EndedAt:     c.EndedAt,
		TimeBaseDomain: V1Domains.TimeBaseDomain{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		},
	}
}

func FromDiscHisotryV1Domain(c *V1Domains.DiscHistoryDomain) DiscountHistories {
	return DiscountHistories{
		ID:          c.ID,
		ProductID:   c.ProductID,
		Price:       c.Price,
		SaveAmount:  c.SaveAmount,
		SavePercent: c.SavePercent,
		Duration:    c.Duration,
		StartedAt:   c.StartedAt,
		EndedAt:     c.EndedAt,
		TimeBase: TimeBase{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		},
	}
}

func ToDiscHisotryV1Domain(c *[]DiscountHistories) []V1Domains.DiscHistoryDomain {
	var result []V1Domains.DiscHistoryDomain

	for _, val := range *c {
		result = append(result, val.ToV1Domain())
	}

	return result
}
