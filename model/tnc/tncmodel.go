package tnc

import (
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TncModel = (*customTncModel)(nil)

type (
	// TncModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTncModel.
	TncModel interface {
		tncModel
		withSession(session sqlx.Session) TncModel
	}

	customTncModel struct {
		*defaultTncModel
	}
)

// NewTncModel returns a model for the database table.
func NewTncModel(conn sqlx.SqlConn) TncModel {
	return &customTncModel{
		defaultTncModel: newTncModel(conn),
	}
}

func (m *customTncModel) withSession(session sqlx.Session) TncModel {
	return NewTncModel(sqlx.NewSqlConnFromSession(session))
}
