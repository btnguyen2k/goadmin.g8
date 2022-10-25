package myapp

// @availabble since template-r3

import (
	"fmt"
	"strings"
	"time"

	prom "github.com/btnguyen2k/prom/sql"
	_ "github.com/go-sql-driver/mysql"
)

func newMysqlConnection(url string, loc *time.Location) *prom.SqlConnect {
	return newSqlConnection("mysql", url, prom.FlavorMySql, loc)
}

/*----------------------------------------------------------------------*/

const (
	mysqlTableGroup = namespace + "_group"
)

var (
	mysqlColNamesAndTypesGroup = []string{"%s VARCHAR(64)", "%s VARCHAR(255)"}
)

func mysqlInitTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sqlStm := "CREATE TABLE IF NOT EXISTS %s (" + strings.Join(mysqlColNamesAndTypesGroup, ",") + ",PRIMARY KEY (%s))"
	sqlStm = fmt.Sprintf(sqlStm, tableName, sqlColGroupId, sqlColGroupName, sqlColGroupId)
	_, err := sqlc.GetDB().Exec(sqlStm)
	if err != nil {
		panic(err)
	}
}

func newGroupDaoMysql(sqlc *prom.SqlConnect, tableName string) GroupDao {
	return newGroupDaoSql(sqlc, tableName)
}

/*----------------------------------------------------------------------*/

const (
	mysqlTableUser = namespace + "_user"
)

var (
	mysqlColNamesAndTypesUser = []string{"%s VARCHAR(64)", "%s VARCHAR(64)", "%s VARCHAR(64)", "%s VARCHAR(64)"}
)

func mysqlInitTableUser(sqlc *prom.SqlConnect, tableName string) {
	sqlStm := "CREATE TABLE IF NOT EXISTS %s (" + strings.Join(mysqlColNamesAndTypesUser, ",") + ",PRIMARY KEY (%s))"
	sqlStm = fmt.Sprintf(sqlStm, tableName, sqlColUserUsername, sqlColUserPassword, sqlColUserName, sqlColUserGroupId, sqlColUserUsername)
	_, err := sqlc.GetDB().Exec(sqlStm)
	if err != nil {
		panic(err)
	}
}

func newUserDaoMysql(sqlc *prom.SqlConnect, tableName string) UserDao {
	return newUserDaoSql(sqlc, tableName)
}
