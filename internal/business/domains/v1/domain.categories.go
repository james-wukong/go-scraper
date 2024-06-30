package v1

import (
	"context"
)

type CategoryDomain struct {
	Id       int
	ParentId *int
	Name     string
	Level    uint
	Url      string
	Platform uint
}

type CategoryUCInterface interface {
	// used by scrapper
	Store(inDom *CategoryDomain) (outDom CategoryDomain, statusCode int, err error)
	// used by handler
	GetByName(ctx context.Context, name string) (outDom CategoryDomain, statusCode int, err error)
}

type CategoryRepoInterface interface {
	// used by scrapper: get categoryId by it's name and platform
	GetByNamePlatform(name string, platform uint) (outDom CategoryDomain, err error)
	// InsertCategory(inDom *records.Categories) (lastInsertId int, err error)
	GetByURLPlatform(url string, platform uint) (outDom CategoryDomain, err error)
	// used by scrapper
	Upsert(inDom *CategoryDomain) (lastInsertId int, err error)
	// used by handler
	// GetByName(name string) (outDom CategoryDomain, err error)
}
