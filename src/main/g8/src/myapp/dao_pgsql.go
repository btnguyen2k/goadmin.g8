package myapp

// @availabble since template-v0.4.r2

import (
	"fmt"
	"strings"
	"time"

	prom "github.com/btnguyen2k/prom/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func newPgsqlConnection(url string, loc *time.Location) *prom.SqlConnect {
	return newSqlConnection("pgx", url, prom.FlavorPgSql, loc)
}

/*----------------------------------------------------------------------*/

const (
	pgsqlTableGroup = namespace + "_group"
)

var (
	pgsqlColNamesAndTypesGroup = []string{"%s VARCHAR(64)", "%s VARCHAR(255)"}
)

func pgsqlInitTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sqlStm := "CREATE TABLE IF NOT EXISTS %s (" + strings.Join(pgsqlColNamesAndTypesGroup, ",") + ",PRIMARY KEY (%s))"
	sqlStm = fmt.Sprintf(sqlStm, tableName, sqlColGroupId, sqlColGroupName, sqlColGroupId)
	_, err := sqlc.GetDB().Exec(sqlStm)
	if err != nil {
		panic(err)
	}
}

func newGroupDaoPgsql(sqlc *prom.SqlConnect, tableName string) GroupDao {
	return newGroupDaoSql(sqlc, tableName)
}

/*----------------------------------------------------------------------*/

const (
	pgsqlTableUser = namespace + "_user"
)

var (
	pgsqlColNamesAndTypesUser = []string{"%s VARCHAR(64)", "%s VARCHAR(64)", "%s VARCHAR(64)", "%s VARCHAR(64)"}
)

func pgsqlInitTableUser(sqlc *prom.SqlConnect, tableName string) {
	sqlStm := "CREATE TABLE IF NOT EXISTS %s (" + strings.Join(pgsqlColNamesAndTypesUser, ",") + ",PRIMARY KEY (%s))"
	sqlStm = fmt.Sprintf(sqlStm, tableName, sqlColUserUsername, sqlColUserPassword, sqlColUserName, sqlColUserGroupId, sqlColUserUsername)
	_, err := sqlc.GetDB().Exec(sqlStm)
	if err != nil {
		panic(err)
	}
}

func newUserDaoPgsql(sqlc *prom.SqlConnect, tableName string) UserDao {
	return newUserDaoSql(sqlc, tableName)
}
