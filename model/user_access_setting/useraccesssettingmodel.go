package user_access_setting

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserAccessSettingModel = (*customUserAccessSettingModel)(nil)

type (
	// UserAccessSettingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAccessSettingModel.
	UserAccessSettingModel interface {
		userAccessSettingModel
		withSession(session sqlx.Session) UserAccessSettingModel
	}

	customUserAccessSettingModel struct {
		*defaultUserAccessSettingModel
	}
)

// NewUserAccessSettingModel returns a model for the database table.
func NewUserAccessSettingModel(conn sqlx.SqlConn) UserAccessSettingModel {
	return &customUserAccessSettingModel{
		defaultUserAccessSettingModel: newUserAccessSettingModel(conn),
	}
}

func (m *customUserAccessSettingModel) withSession(session sqlx.Session) UserAccessSettingModel {
	return NewUserAccessSettingModel(sqlx.NewSqlConnFromSession(session))
}
