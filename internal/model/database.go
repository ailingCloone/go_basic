// model/database.go

package model

import (
	"nrs_customer_module_backend/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/conf"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func InitializeDatabase() (sqlx.SqlConn, error) {
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	connection := c.DevDBStaging.Source

	if c.Environment == "staging" {
		connection = c.DBStaging.Source
	} else if c.Environment == "prod" {
		connection = c.DBProd.Source
	}

	db := sqlx.NewSqlConn("mysql", connection)
	return db, nil

}
