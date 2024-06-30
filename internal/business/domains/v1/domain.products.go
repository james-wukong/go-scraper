package v1

import (
	"context"
)

type ProductDomain struct {
	Id         string
	CategoryId int
	Name       string
	Sku        string
	Model      string
	ProdId     string
	Price      float32
	Source     int
	UrlLink    string
	ImageLink  string
	Detail     *DetailBaseDomain
	Spec       *SpecBaseDomain
	TimeBaseDomain
}

type ProductUCInterface interface {
	// used by handler
	Store(inDom *ProductDomain) (outDom ProductDomain, statusCode int, err error)
	// used by handler
	GetByName(ctx context.Context, name string) (outDom ProductDomain, statusCode int, err error)
}

type ProductRepoInterface interface {
	// used by scrapper: get product id by it's name and sku
	GetByNameSku(name string, sku string) (outDom ProductDomain, err error)
	// used by scrapper
	UpsertProduct(inDom *ProductDomain) (lastInsertId string, err error)
	// used by handler
	// GetByName(name string) (outDom CategoryDomain, err error)
}
