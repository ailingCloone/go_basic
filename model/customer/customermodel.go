package customer

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CustomerModel = (*customCustomerModel)(nil)

type (
	// CustomerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerModel.
	CustomerModel interface {
		customerModel
		withSession(session sqlx.Session) CustomerModel
	}

	customCustomerModel struct {
		*defaultCustomerModel
	}
)

// NewCustomerModel returns a model for the database table.
func NewCustomerModel(conn sqlx.SqlConn) CustomerModel {
	return &customCustomerModel{
		defaultCustomerModel: newCustomerModel(conn),
	}
}

func (m *customCustomerModel) withSession(session sqlx.Session) CustomerModel {
	return NewCustomerModel(sqlx.NewSqlConnFromSession(session))
}
