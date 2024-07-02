package v1

import (
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
	"github.com/james-wukong/go-app/internal/datasources/records"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NewDiscHistoryRepo(conn *sqlx.DB) V1Domains.DiscHistoryRepoInterface {
	return &psqlRepoConn{
		conn: conn,
	}
}

func (r *psqlRepoConn) GetHistoryByPID(pid string) (outDom V1Domains.DiscHistoryDomain, err error) {
	history := records.DiscountHistories{}
	err = r.conn.Get(&history, `SELECT * FROM discount_histories WHERE product_id=$1 ORDER BY id DESC LIMIT 1;`, history)
	if err != nil {
		logger.Debug("avg ratings query error: "+err.Error(), logrus.Fields{"product": history, "product_id": pid})
	}
	// return category.ToV1Domain(), err
	return history.ToV1Domain(), err
}

func (r *psqlRepoConn) SaveDiscHistory(inDom *V1Domains.DiscHistoryDomain) (lastInsertId uint, err error) {
	history := records.FromDiscHisotryV1Domain(inDom)
	var qryFind string = `SELECT id FROM discount_histories where product_id=$1 AND ended_at=$2 AND duration=$3 Limit 1;`
	var qryInsert string = `INSERT INTO discount_histories (product_id, price, save_amount, save_percent, duration, started_at, ended_at, created_at)
		VALUES (:product_id, :price, :save_amount, :save_percent, :duration, :started_at, :ended_at, :created_at)
		RETURNING id;`
	// search for pre-existing
	err = r.conn.Get(&history, qryFind, inDom.ProductID, inDom.EndedAt, inDom.Duration)
	if err == nil {
		// found previous record
		return history.ID, nil
	} else {
		logger.Debug("error getting discount history: ", logrus.Fields{"err": err})
		// insert history record
		// TODO Test the next query -> not working
		// err = r.conn.QueryRowx(qryInsert, productJSON).Scan(&lastInsertId)
		// if _, err = r.conn.NamedQuery(qryInsert, history); err != nil {
		// 	logger.Debug("error inserting discount history: ", logrus.Fields{"err": err})
		// 	return 0, err
		// }
		stmt, err := r.conn.PrepareNamed(qryInsert)
		if err != nil {
			logger.Debug("discount preparenamed insert", logrus.Fields{"err": err})
			return 0, err
		}
		defer stmt.Close()

		if err = stmt.Get(&lastInsertId, history); err != nil {
			logger.Debug("discount preparenamed get", logrus.Fields{"err": err})
			return 0, err
		}
		return lastInsertId, nil
	}

}
