package email_template

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmailTemplateModel = (*customEmailTemplateModel)(nil)

type (
	// EmailTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailTemplateModel.
	EmailTemplateModel interface {
		emailTemplateModel
		withSession(session sqlx.Session) EmailTemplateModel
	}

	customEmailTemplateModel struct {
		*defaultEmailTemplateModel
	}
)

// NewEmailTemplateModel returns a model for the database table.
func NewEmailTemplateModel(conn sqlx.SqlConn) EmailTemplateModel {
	return &customEmailTemplateModel{
		defaultEmailTemplateModel: newEmailTemplateModel(conn),
	}
}

func (m *customEmailTemplateModel) withSession(session sqlx.Session) EmailTemplateModel {
	return NewEmailTemplateModel(sqlx.NewSqlConnFromSession(session))
}
