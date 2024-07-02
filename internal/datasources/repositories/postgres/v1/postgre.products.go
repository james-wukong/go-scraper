package v1

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
	"github.com/james-wukong/go-app/internal/datasources/records"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NewProductRepo(conn *sqlx.DB) V1Domains.ProductRepoInterface {
	return &psqlRepoConn{
		conn: conn,
	}
}

func (r *psqlRepoConn) GetByNameSku(name string, sku string) (outDom V1Domains.ProductDomain, err error) {
	product := records.Products{}
	err = r.conn.Get(&product, `SELECT * FROM products WHERE name=$1 AND sku=$2 LIMIT 1;`, name, sku)
	if err != nil {
		logger.Debug("product query error: ", logrus.Fields{"err": err, "product": product, "name": name, "sku": sku})
	}
	// return category.ToV1Domain(), err
	return product.ToV1Domain(), err
}

func (r *psqlRepoConn) UpsertProduct(inDom *V1Domains.ProductDomain) (lastInsertId string, err error) {
	product := records.FromProductV1Domain(inDom)
	var qryGet string = `SELECT id FROM products WHERE name=$1 AND sku=$2 LIMIT 1;`

	if err = r.conn.Get(&product, qryGet, product.Name, product.Sku); err != nil {
		// can't find a product, then insert new one
		logger.Debug("product query error: ", logrus.Fields{"error": err, "product": product})
		var qryInsert string = `INSERT INTO products (id, category_id, name, sku, model, prod_id, brand, price, source, url_link, img_link, detail, specification, created_at)
		VALUES (uuid_generate_v4(), :category_id, :name, :sku, :model, :prod_id, :brand, :price, :source, :url_link, :img_link, :detail, :specification, :created_at)
		RETURNING id;`
		productJSON := product.ToProductJson()

		stmt, err := r.conn.PrepareNamed(qryInsert)
		if err != nil {
			logger.Debug("product preparenamed insert", logrus.Fields{"err": err})
			return "", err
		}
		defer stmt.Close()

		if err = stmt.Get(&lastInsertId, productJSON); err != nil {
			logger.Debug("product preparenamed get", logrus.Fields{"err": err})
			return "", err
		}
		logger.Debug("last inserted product id", logrus.Fields{"lastInsertId": lastInsertId, "product": product})
		return lastInsertId, nil
	} else {
		logger.Debug("last found product id", logrus.Fields{"found id": product.Id})
		return product.Id, nil
	}
}
