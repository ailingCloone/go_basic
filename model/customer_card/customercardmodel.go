package customer_card

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CustomerCardModel = (*customCustomerCardModel)(nil)

type (
	// CustomerCardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerCardModel.
	CustomerCardModel interface {
		customerCardModel
		withSession(session sqlx.Session) CustomerCardModel
	}

	customCustomerCardModel struct {
		*defaultCustomerCardModel
	}
)

// NewCustomerCardModel returns a model for the database table.
func NewCustomerCardModel(conn sqlx.SqlConn) CustomerCardModel {
	return &customCustomerCardModel{
		defaultCustomerCardModel: newCustomerCardModel(conn),
	}
}

func (m *customCustomerCardModel) withSession(session sqlx.Session) CustomerCardModel {
	return NewCustomerCardModel(sqlx.NewSqlConnFromSession(session))
}
