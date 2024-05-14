// Code generated by goctl. DO NOT EDIT.

package tenant

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"nrs_customer_module_backend/internal/global"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tenantFieldNames          = builder.RawFieldNames(&Tenant{})
	tenantRows                = strings.Join(tenantFieldNames, ",")
	tenantRowsExpectAutoSet   = strings.Join(stringx.Remove(tenantFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tenantRowsWithPlaceHolder = strings.Join(stringx.Remove(tenantFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tenantModel interface {
		Insert(ctx context.Context, data *Tenant) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Tenant, error)
		Update(ctx context.Context, data *Tenant) error
		Delete(ctx context.Context, id int64) error
		FindTenantWithModule(ctx context.Context, client_id, client_secret string,module int) (*FindTenantWithModuleDB, error)

	}

	defaultTenantModel struct {
		conn  sqlx.SqlConn
		table string
		tenantSubscriptionTable string
	}

	Tenant struct {
		Id           int64     `db:"id"`
		Guid         string    `db:"guid"`
		ClientId     string    `db:"client_id"`
		ClientSecret string    `db:"client_secret"`
		Title        string    `db:"title"`
		Updated      time.Time `db:"updated"`
		Created      time.Time `db:"created"`
		Active       int64     `db:"active"` // 0-Inactive, 1- Active
	}

	FindTenantWithModuleDB struct {
		Id           int64     `db:"id"`
		RoleId       int64     `db:"role_id"`
	}
)

func newTenantModel(conn sqlx.SqlConn) *defaultTenantModel {
	return &defaultTenantModel{
		conn:  conn,
		table: "tenant t",
		tenantSubscriptionTable: "tenant_subscription ts",

	}
}

func (m *defaultTenantModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTenantModel) FindOne(ctx context.Context, id int64) (*Tenant, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tenantRows, m.table)
	var resp Tenant
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

func (m *defaultTenantModel) Insert(ctx context.Context, data *Tenant) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, tenantRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Guid, data.ClientId, data.ClientSecret, data.Title, data.Updated, data.Created, data.Active)
	return ret, err
}

func (m *defaultTenantModel) Update(ctx context.Context, data *Tenant) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tenantRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Guid, data.ClientId, data.ClientSecret, data.Title, data.Updated, data.Created, data.Active, data.Id)
	return err
}

func (m *defaultTenantModel) tableName() string {
	return m.table
}


func (m *defaultTenantModel) FindTenantWithModule(ctx context.Context, client_id, client_secret string,module int) (*FindTenantWithModuleDB, error) {
		currentTime:= global.GetCurrentTime()
	
		client_secret = global.GenerateSha256(client_secret)
	
		query := fmt.Sprintf("SELECT id,role_id FROM %s INNER JOIN %s ON t.id = ts.tenant_id AND ts.active = 1 AND ts.main_module = ? AND ts.date_from <= ? AND ts.date_to >= ? WHERE t.active = 1 AND t.client_id = ? AND t.client_secret = ?  LIMIT 1",  m.table,m.tenantSubscriptionTable)
		var resp FindTenantWithModuleDB
		err := m.conn.QueryRowCtx(ctx, &resp, query,module, currentTime,currentTime, client_id,client_secret)
		switch err {
		case nil:
			return &resp, nil
		case sqlx.ErrNotFound:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}