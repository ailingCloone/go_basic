package gs_country_code

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GsCountryCodeModel = (*customGsCountryCodeModel)(nil)

type (
	// GsCountryCodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGsCountryCodeModel.
	GsCountryCodeModel interface {
		gsCountryCodeModel
		withSession(session sqlx.Session) GsCountryCodeModel
	}

	customGsCountryCodeModel struct {
		*defaultGsCountryCodeModel
	}
)

// NewGsCountryCodeModel returns a model for the database table.
func NewGsCountryCodeModel(conn sqlx.SqlConn) GsCountryCodeModel {
	return &customGsCountryCodeModel{
		defaultGsCountryCodeModel: newGsCountryCodeModel(conn),
	}
}

func (m *customGsCountryCodeModel) withSession(session sqlx.Session) GsCountryCodeModel {
	return NewGsCountryCodeModel(sqlx.NewSqlConnFromSession(session))
}
