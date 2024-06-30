package v1

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
	"github.com/james-wukong/go-app/internal/datasources/records"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NewCategoryRepo(conn *sqlx.DB) V1Domains.CategoryRepoInterface {
	return &psqlRepoConn{
		conn: conn,
	}
}

func (r *psqlRepoConn) GetByNamePlatform(name string, platform uint) (outDom V1Domains.CategoryDomain, err error) {
	category := records.Categories{}
	err = r.conn.Get(&category, "SELECT * FROM categories WHERE name=$1 AND platform=$2 LIMIT 1;", name, platform)

	return category.ToV1Domain(), err
}

func (r *psqlRepoConn) GetByURLPlatform(url string, platform uint) (outDom V1Domains.CategoryDomain, err error) {
	category := records.Categories{}
	// logger.Debug("executing GetByURLPlatform", logrus.Fields{"url": url})
	err = r.conn.Get(&category, "SELECT * FROM categories WHERE url=$1 AND platform=$2 LIMIT 1;", url, platform)

	return category.ToV1Domain(), err
}

func (r *psqlRepoConn) insertCategory(inDom *records.Categories) (lastInsertId int, err error) {
	var query string = `INSERT INTO categories (parent_id, name, level, url, platform) 
	VALUES (:parent_id, :name, :level, :url, :platform) 
	RETURNING id`
	var qryGet string = `SELECT * FROM categories WHERE parent_id=$1 AND name=$2 AND platform=$3 LIMIT 1;`
	var qryGetNil string = `SELECT * FROM categories WHERE parent_id IS NULL AND name=$1 AND platform=$2 LIMIT 1;`

	// Additional debug log to ensure the query and input data are correct before executing the query
	logger.Debug("executing insert query", logrus.Fields{
		"query": query,
		"data":  inDom,
		"parent_id": func() interface{} {
			if inDom.ParentId != nil {
				return *inDom.ParentId
			}
			return nil
		}(),
	})

	// _, err = r.conn.NamedExec(`INSERT INTO categories (parent_id, name, level, url, platform)
	// VALUES (:parent_id, :name, :level, :url, :platform)`, category)
	if _, err = r.conn.NamedExec(query, inDom); err != nil {
		return 0, err
	}
	if inDom.ParentId == nil {
		if err = r.conn.Get(inDom, qryGetNil, inDom.Name, inDom.Platform); err != nil {
			return 0, err
		}
	} else {
		if err = r.conn.Get(inDom, qryGet, inDom.ParentId, inDom.Name, inDom.Platform); err != nil {
			return 0, err
		}
	}
	return inDom.Id, nil
}

func (r *psqlRepoConn) Upsert(inDom *V1Domains.CategoryDomain) (lastInsertId int, err error) {
	category := records.FromCategoryV1Domain(inDom)
	var qryGet string = `SELECT id FROM categories WHERE parent_id=$1 AND name=$2 AND platform=$3 LIMIT 1;`
	var qryGetNil string = `SELECT id FROM categories WHERE parent_id IS NULL AND name=$1 AND platform=$2 LIMIT 1;`
	if category.ParentId == nil {
		if err = r.conn.Get(&category, qryGetNil, category.Name, category.Platform); err != nil {
			lastInsertId, err = r.insertCategory(&category)

			return lastInsertId, err
		}
	} else {
		if err = r.conn.Get(&category, qryGet, *category.ParentId, category.Name, category.Platform); err != nil {
			lastInsertId, err = r.insertCategory(&category)

			return lastInsertId, err
		}
	}

	return category.Id, nil
}
