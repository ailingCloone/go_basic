package activity_log

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ActivityLogModel = (*customActivityLogModel)(nil)

type (
	// ActivityLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customActivityLogModel.
	ActivityLogModel interface {
		activityLogModel
		withSession(session sqlx.Session) ActivityLogModel
	}

	customActivityLogModel struct {
		*defaultActivityLogModel
	}
)

// NewActivityLogModel returns a model for the database table.
func NewActivityLogModel(conn sqlx.SqlConn) ActivityLogModel {
	return &customActivityLogModel{
		defaultActivityLogModel: newActivityLogModel(conn),
	}
}

func (m *customActivityLogModel) withSession(session sqlx.Session) ActivityLogModel {
	return NewActivityLogModel(sqlx.NewSqlConnFromSession(session))
}
