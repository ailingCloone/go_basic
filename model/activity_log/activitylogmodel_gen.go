// Code generated by goctl. DO NOT EDIT.

package activity_log

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	// "time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	activityLogFieldNames          = builder.RawFieldNames(&ActivityLog{})
	activityLogRows                = strings.Join(activityLogFieldNames, ",")
	activityLogRowsExpectAutoSet   = strings.Join(stringx.Remove(activityLogFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	activityLogRowsWithPlaceHolder = strings.Join(stringx.Remove(activityLogFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	activityLogModel interface {
		Insert(ctx context.Context, data *ActivityLog) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*ActivityLog, error)
		Update(ctx context.Context, data *ActivityLog) error
		Delete(ctx context.Context, id int64) error
	}

	defaultActivityLogModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ActivityLog struct {
		Id         int64     `db:"id"`
		StaffId    int64     `db:"staff_id"`
		CustomerId int64     `db:"customer_id"`
		ReferTable string    `db:"refer_table"`
		Action     string    `db:"action"`
		Changes    string    `db:"changes"` // Ex: {"colum_name":[{"before":"","after":"",}]}
		Created    string    `db:"created"`
	}

	// Create a structure to represent the changes
	 Change struct {
		Before string `json:"before"`
		After  string `json:"after"`
	}

)

func newActivityLogModel(conn sqlx.SqlConn) *defaultActivityLogModel {
	return &defaultActivityLogModel{
		conn:  conn,
		table: "`activity_log`",
	}
}

func (m *defaultActivityLogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultActivityLogModel) FindOne(ctx context.Context, id int64) (*ActivityLog, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", activityLogRows, m.table)
	var resp ActivityLog
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultActivityLogModel) Insert(ctx context.Context, data *ActivityLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, activityLogRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.StaffId, data.CustomerId, data.ReferTable, data.Action, data.Changes, data.Created)
	return ret, err
}

func (m *defaultActivityLogModel) Update(ctx context.Context, data *ActivityLog) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, activityLogRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.StaffId, data.CustomerId, data.ReferTable, data.Action, data.Changes, data.Created, data.Id)
	return err
}

func (m *defaultActivityLogModel) tableName() string {
	return m.table
}
