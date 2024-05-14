// Code generated by goctl. DO NOT EDIT.

package staff

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	staffFieldNames          = builder.RawFieldNames(&Staff{})
	staffRows                = strings.Join(staffFieldNames, ",")
	staffRowsExpectAutoSet   = strings.Join(stringx.Remove(staffFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	staffRowsWithPlaceHolder = strings.Join(stringx.Remove(staffFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	staffModel interface {
		Insert(ctx context.Context, data *Staff) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Staff, error)
		Update(ctx context.Context, data *Staff) error
		Delete(ctx context.Context, id int64) error
		FindOneGuid(ctx context.Context, guid string) (*Staff, error) 
		FindOneEmail(ctx context.Context, email string) (*Staff, error) 
		FindOneContact(ctx context.Context, contact string) (*Staff, error)
		FindOneIC(ctx context.Context, ic string) (*Staff, error)
		UpdatePassword(ctx context.Context, data *Staff) error	
	}

	defaultStaffModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Staff struct {
		Id          int64     `db:"id"`
		Guid        string    `db:"guid"`
		Code        string    `db:"code"`
		Password    string    `db:"password"`
		Outlet      int64     `db:"outlet"`
		Role        int64     `db:"role"`
		Name        string    `db:"name"`
		CountryCode string    `db:"country_code"`
		Icno        string    `db:"icno"`
		Contact     string    `db:"contact"`
		Email       string    `db:"email"`
		Address     string    `db:"address"`
		WorkingFrom time.Time `db:"working_from"`
		ConfirmDate time.Time `db:"confirm_date"`
		WorkingTo   time.Time `db:"working_to"`
		Updated     time.Time `db:"updated"`
		Created     time.Time `db:"created"`
		Active      int64     `db:"active"` // 0:Inactive, 1-Active
	}

	StaffInfo struct {
		Guid      string `json:"guid"`
		ImageUrl  string `json:"image_url"`
		WebImageUrl  string `json:"web_image_url"`
		Name      string `json:"name"`
		Contact   string `json:"contact"`
		Email     string `json:"email"`
		Icno      string `json:"icno"`
		StaffCode string `json:"staff_code"`
	}
	
)

func newStaffModel(conn sqlx.SqlConn) *defaultStaffModel {
	return &defaultStaffModel{
		conn:  conn,
		table: "`staff`",
	}
}

func (m *defaultStaffModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultStaffModel) FindOne(ctx context.Context, id int64) (*Staff, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", staffRows, m.table)
	var resp Staff
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

func (m *defaultStaffModel) Insert(ctx context.Context, data *Staff) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, staffRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Guid, data.Code, data.Password, data.Outlet, data.Role, data.Name, data.CountryCode, data.Icno, data.Contact, data.Email, data.Address, data.WorkingFrom, data.ConfirmDate, data.WorkingTo, data.Updated, data.Created, data.Active)
	return ret, err
}

func (m *defaultStaffModel) Update(ctx context.Context, data *Staff) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, staffRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Guid, data.Code, data.Password, data.Outlet, data.Role, data.Name, data.CountryCode, data.Icno, data.Contact, data.Email, data.Address, data.WorkingFrom, data.ConfirmDate, data.WorkingTo, data.Updated, data.Created, data.Active, data.Id)
	return err
}

func (m *defaultStaffModel) tableName() string {
	return m.table
}

func (m *defaultStaffModel) FindOneGuid(ctx context.Context, guid string) (*Staff, error) {
	query := fmt.Sprintf("select %s from %s where active  = 1 and `guid` = ? limit 1", staffRows, m.table)
	var resp Staff
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

func (m *defaultStaffModel) FindOneEmail(ctx context.Context, email string) (*Staff, error) {
    query := fmt.Sprintf("select %s from %s where `email` = ? AND `active` = 1 limit 1", staffRows, m.table)
    var resp Staff
    err := m.conn.QueryRowCtx(ctx, &resp, query, email)
    switch err {
    case nil:
        return &resp, nil
    case sqlx.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *defaultStaffModel) UpdatePassword(ctx context.Context, data *Staff) error {
	query := fmt.Sprintf("UPDATE %s SET `password` = ?, `updated` = ? WHERE `email` = ? AND `active` = 1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Password,data.Updated, data.Email)
	return err
}

func (m *defaultStaffModel) FindOneContact(ctx context.Context, contact string) (*Staff, error) {
    query := fmt.Sprintf("select %s from %s where `contact` = ? AND `active` = 1 limit 1", staffRows, m.table)
    var resp Staff
    err := m.conn.QueryRowCtx(ctx, &resp, query, contact)
    switch err {
    case nil:
        return &resp, nil
    case sqlx.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *defaultStaffModel) FindOneIC(ctx context.Context, ic string) (*Staff, error) {
    query := fmt.Sprintf("select %s from %s where `icno` = ? AND `active` = 1 limit 1", staffRows, m.table)
    var resp Staff
    err := m.conn.QueryRowCtx(ctx, &resp, query, ic)
    switch err {
    case nil:
        return &resp, nil
    case sqlx.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}
