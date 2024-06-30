package records

import (
	"encoding/json"

	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
)

func (c *Products) ToV1Domain() V1Domains.ProductDomain {
	return V1Domains.ProductDomain{
		Id:         c.Id,
		CategoryId: c.CategoryId,
		Name:       c.Name,
		Sku:        c.Sku,
		Model:      c.Model,
		ProdId:     c.ProdId,
		Price:      c.Price,
		Source:     c.Source,
		UrlLink:    c.UrlLink,
		ImageLink:  c.ImageLink,
		Detail:     &V1Domains.DetailBaseDomain{Detail: c.Detail.Detail},
		Spec:       &V1Domains.SpecBaseDomain{Spec: c.Spec.Spec},
		TimeBaseDomain: V1Domains.TimeBaseDomain{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		},
	}
}

func FromProductV1Domain(c *V1Domains.ProductDomain) Products {
	return Products{
		Id:         c.Id,
		CategoryId: c.CategoryId,
		Name:       c.Name,
		Sku:        c.Sku,
		Model:      c.Model,
		ProdId:     c.ProdId,
		Price:      c.Price,
		Source:     c.Source,
		UrlLink:    c.UrlLink,
		ImageLink:  c.ImageLink,
		Detail:     &DetailBase{Detail: c.Detail.Detail},
		Spec:       &SpecBase{Spec: c.Spec.Spec},
		TimeBase: TimeBase{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		},
	}
}

func (c *Products) ToProductJson() map[string]any {
	var detailJSON, specJSON []byte
	var err error
	if detailJSON, err = json.Marshal(c.Detail); err != nil {
		detailJSON = []byte{}
	}
	if specJSON, err = json.Marshal(c.Spec); err != nil {
		specJSON = []byte{}
	}

	return map[string]interface{}{
		"category_id":   c.CategoryId,
		"name":          c.Name,
		"sku":           c.Sku,
		"model":         c.Model,
		"prod_id":       c.ProdId,
		"price":         c.Price,
		"source":        c.Source,
		"url_link":      c.UrlLink,
		"img_link":      c.ImageLink,
		"detail":        detailJSON,
		"specification": specJSON,
		"created_at":    c.CreatedAt,
	}
}

func ToProductsV1Domain(c *[]Products) []V1Domains.ProductDomain {
	var result []V1Domains.ProductDomain

	for _, val := range *c {
		result = append(result, val.ToV1Domain())
	}

	return result
}
