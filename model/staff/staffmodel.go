package staff

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StaffModel = (*customStaffModel)(nil)

type (
	// StaffModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStaffModel.
	StaffModel interface {
		staffModel
		withSession(session sqlx.Session) StaffModel
	}

	customStaffModel struct {
		*defaultStaffModel
	}
)

// NewStaffModel returns a model for the database table.
func NewStaffModel(conn sqlx.SqlConn) StaffModel {
	return &customStaffModel{
		defaultStaffModel: newStaffModel(conn),
	}
}

func (m *customStaffModel) withSession(session sqlx.Session) StaffModel {
	return NewStaffModel(sqlx.NewSqlConnFromSession(session))
}
