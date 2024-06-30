package v1

import "github.com/jmoiron/sqlx"

type psqlRepoConn struct {
	conn *sqlx.DB
}
