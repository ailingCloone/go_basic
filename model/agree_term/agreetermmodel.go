package agree_term

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AgreeTermModel = (*customAgreeTermModel)(nil)

type (
	// AgreeTermModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAgreeTermModel.
	AgreeTermModel interface {
		agreeTermModel
		withSession(session sqlx.Session) AgreeTermModel
	}

	customAgreeTermModel struct {
		*defaultAgreeTermModel
	}
)

// NewAgreeTermModel returns a model for the database table.
func NewAgreeTermModel(conn sqlx.SqlConn) AgreeTermModel {
	return &customAgreeTermModel{
		defaultAgreeTermModel: newAgreeTermModel(conn),
	}
}

func (m *customAgreeTermModel) withSession(session sqlx.Session) AgreeTermModel {
	return NewAgreeTermModel(sqlx.NewSqlConnFromSession(session))
}
