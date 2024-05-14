package homepage

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomepageModel = (*customHomepageModel)(nil)

type (
	// HomepageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomepageModel.
	HomepageModel interface {
		homepageModel
		withSession(session sqlx.Session) HomepageModel
	}

	customHomepageModel struct {
		*defaultHomepageModel
	}
)

// NewHomepageModel returns a model for the database table.
func NewHomepageModel(conn sqlx.SqlConn) HomepageModel {
	return &customHomepageModel{
		defaultHomepageModel: newHomepageModel(conn),
	}
}

func (m *customHomepageModel) withSession(session sqlx.Session) HomepageModel {
	return NewHomepageModel(sqlx.NewSqlConnFromSession(session))
}
