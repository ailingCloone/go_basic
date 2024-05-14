// Code generated by goctl. DO NOT EDIT.

package homepage

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
	homepageFieldNames          = builder.RawFieldNames(&Homepage{})
	homepageRows                = strings.Join(homepageFieldNames, ",")
	homepageRowsExpectAutoSet   = strings.Join(stringx.Remove(homepageFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	homepageRowsWithPlaceHolder = strings.Join(stringx.Remove(homepageFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	homepageModel interface {
		Insert(ctx context.Context, data *Homepage) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Homepage, error)
		FindOneByGuid(ctx context.Context, guid string) (*Homepage, error)
		FindOneByCategoryId(ctx context.Context, categoryId int64) (*Homepage, error)
		Update(ctx context.Context, data *Homepage) error
		Delete(ctx context.Context, id int64) error
	}

	defaultHomepageModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Homepage struct {
		Id          int64     `db:"id"`
		Guid        string    `db:"guid"`
		ImageUrl    string    `db:"image_url"`
		WebImageUrl string    `db:"web_image_url"`
		Title       string    `db:"title"`
		Action      string    `db:"action"`   // Action  like “./login”
		Priority    int64     `db:"priority"` // The record will order based on this column, 1,2,3
		CategoryId  int64     `db:"category_id"`
		DisplayFrom time.Time `db:"display_from"` // Display from date time
		DisplayTo   time.Time `db:"display_to"`   // Display until date time
		ParentId    int64     `db:"parent_id"`    // refer to sidemenu id: Is submenu of parent
		Updated     time.Time `db:"updated"`
		Created     time.Time `db:"created"`
		Platform    int64     `db:"platform"` // 1- App 2-Portal
		Role        string    `db:"role"`
		Active      int64     `db:"active"` // 0- Inactive , 1- Active
	}
)

func newHomepageModel(conn sqlx.SqlConn) *defaultHomepageModel {
	return &defaultHomepageModel{
		conn:  conn,
		table: "`homepage`",
	}
}

func (m *defaultHomepageModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultHomepageModel) FindOne(ctx context.Context, id int64) (*Homepage, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homepageRows, m.table)
	var resp Homepage
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

func (m *defaultHomepageModel) FindOneByGuid(ctx context.Context, guid string) (*Homepage, error) {
	var resp Homepage
	query := fmt.Sprintf("select %s from %s where `guid` = ? limit 1", homepageRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, guid)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultHomepageModel) FindOneByCategoryId(ctx context.Context, categoryId int64) (*Homepage, error) {
	var resp Homepage
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("select %s from %s where `category_id` = ? and `display_from` <= ? and `display_to` >= ? and `active` = 1 limit 1", homepageRows, m.table)
	err = m.conn.QueryRowCtx(ctx, &resp, query, categoryId, currentTime, currentTime)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultHomepageModel) Insert(ctx context.Context, data *Homepage) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, homepageRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Guid, data.ImageUrl, data.WebImageUrl, data.Title, data.Action, data.Priority, data.CategoryId, data.DisplayFrom, data.DisplayTo, data.ParentId, data.Updated, data.Created, data.Platform, data.Role, data.Active)
	return ret, err
}

func (m *defaultHomepageModel) Update(ctx context.Context, newData *Homepage) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homepageRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Guid, newData.ImageUrl, newData.WebImageUrl, newData.Title, newData.Action, newData.Priority, newData.CategoryId, newData.DisplayFrom, newData.DisplayTo, newData.ParentId, newData.Updated, newData.Created, newData.Platform, newData.Role, newData.Active, newData.Id)
	return err
}

func (m *defaultHomepageModel) tableName() string {
	return m.table
}
