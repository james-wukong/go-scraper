package v1

import (
	"context"

	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
	"github.com/james-wukong/go-app/internal/datasources/records"
	"github.com/jmoiron/sqlx"
)

func NewUserRepository(conn *sqlx.DB) V1Domains.UserRepository {
	return &psqlRepoConn{
		conn: conn,
	}
}

func (r *psqlRepoConn) Store(ctx context.Context, inDom *V1Domains.UserDomain) (err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	_, err = r.conn.NamedQueryContext(ctx, `INSERT INTO users(id, username, email, password, active, role_id, created_at) VALUES (uuid_generate_v4(), :username, :email, :password, false, :role_id, :created_at)`, userRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r *psqlRepoConn) GetByEmail(ctx context.Context, inDom *V1Domains.UserDomain) (outDomain V1Domains.UserDomain, err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	err = r.conn.GetContext(ctx, &userRecord, `SELECT * FROM users WHERE "email" = $1`, userRecord.Email)
	if err != nil {
		return V1Domains.UserDomain{}, err
	}

	return userRecord.ToV1Domain(), nil
}

func (r *psqlRepoConn) ChangeActiveUser(ctx context.Context, inDom *V1Domains.UserDomain) (err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	_, err = r.conn.NamedQueryContext(ctx, `UPDATE users SET active = :active WHERE id = :id`, userRecord)

	return
}
