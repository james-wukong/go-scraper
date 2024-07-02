package v1

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
	"github.com/james-wukong/go-app/internal/datasources/records"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NewAvgRatingRepo(conn *sqlx.DB) V1Domains.AvgRatingRepoInterface {
	return &psqlRepoConn{
		conn: conn,
	}
}

func (r *psqlRepoConn) GetAvgRatingByPID(productID string) (outDom V1Domains.AvgRatingDomain, err error) {
	avgRating := records.AvgRatings{}
	err = r.conn.Get(&avgRating, `SELECT * FROM avg_ratings WHERE product_id=$1 LIMIT 1;`, productID)
	if err != nil {
		logger.Debug("avg ratings query error: "+err.Error(), logrus.Fields{"product": avgRating, "product_id": productID})
		// TODO insert avg rating row
	}
	// return category.ToV1Domain(), err
	return avgRating.ToV1Domain(), err
}

func (r *psqlRepoConn) UpsertAvgRating(inDom *V1Domains.AvgRatingDomain) (lastInsertId uint, err error) {
	avgRating := records.FromAvgRatingV1Domain(inDom)
	// var qryGet string = `SELECT id FROM avg_ratings WHERE product_id=$1 LIMIT 1;`
	var qryDel string = `DELETE FROM avg_ratings WHERE product_id = :product_id;`

	if _, err = r.conn.NamedExec(qryDel, avgRating); err != nil {
		// can't find a avg ratings, then insert new one
		logger.Debug("avg ratings delete error", logrus.Fields{"avgRating": avgRating, "err": err})
	}
	var qryInsert string = `INSERT INTO avg_ratings (product_id, overall, quality, value, total_reviews, star5, star4, star3, star2, star1, created_at)
		VALUES (:product_id, :overall, :quality, :value, :total_reviews, :star5, :star4, :star3, :star2, :star1, :created_at)
		RETURNING id;`
	stmt, err := r.conn.PrepareNamed(qryInsert)
	if err != nil {
		logger.Debug("avg rating preparenamed insert", logrus.Fields{"err": err})
		return 0, err
	}
	defer stmt.Close()

	if err = stmt.Get(&lastInsertId, avgRating); err != nil {
		logger.Debug("avg rating preparenamed get", logrus.Fields{"err": err})
		return 0, err
	}

	return lastInsertId, nil
}
