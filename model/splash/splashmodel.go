package splash

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SplashModel = (*customSplashModel)(nil)

type (
	// SplashModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSplashModel.
	SplashModel interface {
		splashModel
		withSession(session sqlx.Session) SplashModel
	}

	customSplashModel struct {
		*defaultSplashModel
	}
)

// NewSplashModel returns a model for the database table.
func NewSplashModel(conn sqlx.SqlConn) SplashModel {
	return &customSplashModel{
		defaultSplashModel: newSplashModel(conn),
	}
}

func (m *customSplashModel) withSession(session sqlx.Session) SplashModel {
	return NewSplashModel(sqlx.NewSqlConnFromSession(session))
}
