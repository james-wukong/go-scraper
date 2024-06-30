package records

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
)

func (c *Categories) ToV1Domain() V1Domains.CategoryDomain {
	return V1Domains.CategoryDomain{
		Id:       c.Id,
		ParentId: c.ParentId,
		Name:     c.Name,
		Level:    c.Level,
		Url:      c.Url,
		Platform: c.Platform,
	}
}

func FromCategoryV1Domain(c *V1Domains.CategoryDomain) Categories {
	return Categories{
		Id:       c.Id,
		ParentId: c.ParentId,
		Name:     c.Name,
		Level:    c.Level,
		Url:      c.Url,
		Platform: c.Platform,
	}
}

func ToCategoriesV1Domain(c *[]Categories) []V1Domains.CategoryDomain {
	var result []V1Domains.CategoryDomain

	for _, val := range *c {
		result = append(result, val.ToV1Domain())
	}

	return result
}
