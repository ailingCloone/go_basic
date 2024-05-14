package customer_other_details

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CustomerOtherDetailsModel = (*customCustomerOtherDetailsModel)(nil)

type (
	// CustomerOtherDetailsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerOtherDetailsModel.
	CustomerOtherDetailsModel interface {
		customerOtherDetailsModel
		withSession(session sqlx.Session) CustomerOtherDetailsModel
	}

	customCustomerOtherDetailsModel struct {
		*defaultCustomerOtherDetailsModel
	}
)

// NewCustomerOtherDetailsModel returns a model for the database table.
func NewCustomerOtherDetailsModel(conn sqlx.SqlConn) CustomerOtherDetailsModel {
	return &customCustomerOtherDetailsModel{
		defaultCustomerOtherDetailsModel: newCustomerOtherDetailsModel(conn),
	}
}

func (m *customCustomerOtherDetailsModel) withSession(session sqlx.Session) CustomerOtherDetailsModel {
	return NewCustomerOtherDetailsModel(sqlx.NewSqlConnFromSession(session))
}
