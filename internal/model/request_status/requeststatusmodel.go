package request_status

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RequestStatusModel = (*customRequestStatusModel)(nil)

type (
	// RequestStatusModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRequestStatusModel.
	RequestStatusModel interface {
		requestStatusModel
		withSession(session sqlx.Session) RequestStatusModel
	}

	customRequestStatusModel struct {
		*defaultRequestStatusModel
	}
)

// NewRequestStatusModel returns a model for the database table.
func NewRequestStatusModel(conn sqlx.SqlConn) RequestStatusModel {
	return &customRequestStatusModel{
		defaultRequestStatusModel: newRequestStatusModel(conn),
	}
}

func (m *customRequestStatusModel) withSession(session sqlx.Session) RequestStatusModel {
	return NewRequestStatusModel(sqlx.NewSqlConnFromSession(session))
}
