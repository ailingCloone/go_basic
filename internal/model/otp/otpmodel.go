package otp

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OtpModel = (*customOtpModel)(nil)

type (
	// OtpModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOtpModel.
	OtpModel interface {
		otpModel
		withSession(session sqlx.Session) OtpModel
	}

	customOtpModel struct {
		*defaultOtpModel
	}
)

// NewOtpModel returns a model for the database table.
func NewOtpModel(conn sqlx.SqlConn) OtpModel {
	return &customOtpModel{
		defaultOtpModel: newOtpModel(conn),
	}
}

func (m *customOtpModel) withSession(session sqlx.Session) OtpModel {
	return NewOtpModel(sqlx.NewSqlConnFromSession(session))
}
