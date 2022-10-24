package myapp

import (
	"strings"
	"time"

	"github.com/btnguyen2k/prom/sql"
)

const (
	envSqliteDriver = "SQLITE_DRIVER"
	envSqliteUrl    = "SQLITE_URL"
	envMysqlDriver  = "MYSQL_DRIVER"
	envMysqlUrl     = "MYSQL_URL"
	envPgsqlDriver  = "PGSQL_DRIVER"
	envPgsqlUrl     = "PGSQL_URL"
	envCosmosDriver = "COSMOSDB_DRIVER"
	envCosmosUrl    = "COSMOSDB_URL"
)

var (
	testTimeZone = "Asia/_Ho_Chi_Minh"
)

type _m map[string]interface{}

func _newSqlConnect(driver, url, timezone string, flavor sql.DbFlavor) (*sql.SqlConnect, error) {
	driver = strings.Trim(driver, "\"")
	url = strings.Trim(url, "\"")
	if driver == "" || url == "" {
		return nil, nil
	}

	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	sqlc, err := sql.NewSqlConnectWithFlavor(driver, url, 10000, nil, flavor)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}
