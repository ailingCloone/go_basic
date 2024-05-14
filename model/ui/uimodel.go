package ui

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UiModel = (*customUiModel)(nil)

type (
	// UiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUiModel.
	UiModel interface {
		uiModel
		withSession(session sqlx.Session) UiModel
	}

	customUiModel struct {
		*defaultUiModel
	}
)

// NewUiModel returns a model for the database table.
func NewUiModel(conn sqlx.SqlConn) UiModel {
	return &customUiModel{
		defaultUiModel: newUiModel(conn),
	}
}

func (m *customUiModel) withSession(session sqlx.Session) UiModel {
	return NewUiModel(sqlx.NewSqlConnFromSession(session))
}
