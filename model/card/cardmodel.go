package card

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CardModel = (*customCardModel)(nil)

type (
	// CardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCardModel.
	CardModel interface {
		cardModel
		withSession(session sqlx.Session) CardModel
	}

	customCardModel struct {
		*defaultCardModel
	}
)

// NewCardModel returns a model for the database table.
func NewCardModel(conn sqlx.SqlConn) CardModel {
	return &customCardModel{
		defaultCardModel: newCardModel(conn),
	}
}

func (m *customCardModel) withSession(session sqlx.Session) CardModel {
	return NewCardModel(sqlx.NewSqlConnFromSession(session))
}
