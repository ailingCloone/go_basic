package tenant

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TenantModel = (*customTenantModel)(nil)

type (
	// TenantModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTenantModel.
	TenantModel interface {
		tenantModel
		withSession(session sqlx.Session) TenantModel
	}

	customTenantModel struct {
		*defaultTenantModel
	}
)

// NewTenantModel returns a model for the database table.
func NewTenantModel(conn sqlx.SqlConn) TenantModel {
	return &customTenantModel{
		defaultTenantModel: newTenantModel(conn),
	}
}

func (m *customTenantModel) withSession(session sqlx.Session) TenantModel {
	return NewTenantModel(sqlx.NewSqlConnFromSession(session))
}
