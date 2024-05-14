package card_flow

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CardFlowModel = (*customCardFlowModel)(nil)

type (
	// CardFlowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCardFlowModel.
	CardFlowModel interface {
		cardFlowModel
		withSession(session sqlx.Session) CardFlowModel
	}

	customCardFlowModel struct {
		*defaultCardFlowModel
	}
)

// NewCardFlowModel returns a model for the database table.
func NewCardFlowModel(conn sqlx.SqlConn) CardFlowModel {
	return &customCardFlowModel{
		defaultCardFlowModel: newCardFlowModel(conn),
	}
}

func (m *customCardFlowModel) withSession(session sqlx.Session) CardFlowModel {
	return NewCardFlowModel(sqlx.NewSqlConnFromSession(session))
}
