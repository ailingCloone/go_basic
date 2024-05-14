package staff_other_details

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StaffOtherDetailsModel = (*customStaffOtherDetailsModel)(nil)

type (
	// StaffOtherDetailsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStaffOtherDetailsModel.
	StaffOtherDetailsModel interface {
		staffOtherDetailsModel
		withSession(session sqlx.Session) StaffOtherDetailsModel
	}

	customStaffOtherDetailsModel struct {
		*defaultStaffOtherDetailsModel
	}
)

// NewStaffOtherDetailsModel returns a model for the database table.
func NewStaffOtherDetailsModel(conn sqlx.SqlConn) StaffOtherDetailsModel {
	return &customStaffOtherDetailsModel{
		defaultStaffOtherDetailsModel: newStaffOtherDetailsModel(conn),
	}
}

func (m *customStaffOtherDetailsModel) withSession(session sqlx.Session) StaffOtherDetailsModel {
	return NewStaffOtherDetailsModel(sqlx.NewSqlConnFromSession(session))
}
