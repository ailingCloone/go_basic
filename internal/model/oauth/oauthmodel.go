package oauth

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OauthModel = (*customOauthModel)(nil)

type (
	// OauthModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOauthModel.
	OauthModel interface {
		oauthModel
		withSession(session sqlx.Session) OauthModel
	}

	customOauthModel struct {
		*defaultOauthModel
	}
)

// NewOauthModel returns a model for the database table.
func NewOauthModel(conn sqlx.SqlConn) OauthModel {
	return &customOauthModel{
		defaultOauthModel: newOauthModel(conn),
	}
}

func (m *customOauthModel) withSession(session sqlx.Session) OauthModel {
	return NewOauthModel(sqlx.NewSqlConnFromSession(session))
}
