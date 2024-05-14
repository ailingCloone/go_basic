package sidemenu

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SidemenuModel = (*customSidemenuModel)(nil)

type (
	// SidemenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSidemenuModel.
	SidemenuModel interface {
		sidemenuModel
		withSession(session sqlx.Session) SidemenuModel
	}

	customSidemenuModel struct {
		*defaultSidemenuModel
	}
)

// NewSidemenuModel returns a model for the database table.
func NewSidemenuModel(conn sqlx.SqlConn) SidemenuModel {
	return &customSidemenuModel{
		defaultSidemenuModel: newSidemenuModel(conn),
	}
}

func (m *customSidemenuModel) withSession(session sqlx.Session) SidemenuModel {
	return NewSidemenuModel(sqlx.NewSqlConnFromSession(session))
}
