package sms_template

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SmsTemplateModel = (*customSmsTemplateModel)(nil)

type (
	// SmsTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsTemplateModel.
	SmsTemplateModel interface {
		smsTemplateModel
		withSession(session sqlx.Session) SmsTemplateModel
	}

	customSmsTemplateModel struct {
		*defaultSmsTemplateModel
	}
)

// NewSmsTemplateModel returns a model for the database table.
func NewSmsTemplateModel(conn sqlx.SqlConn) SmsTemplateModel {
	return &customSmsTemplateModel{
		defaultSmsTemplateModel: newSmsTemplateModel(conn),
	}
}

func (m *customSmsTemplateModel) withSession(session sqlx.Session) SmsTemplateModel {
	return NewSmsTemplateModel(sqlx.NewSqlConnFromSession(session))
}
