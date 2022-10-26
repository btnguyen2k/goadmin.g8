package myapp

import (
	"fmt"
	"os"
	"strings"
	"time"

	prom "github.com/btnguyen2k/prom/sql"
	_ "github.com/mattn/go-sqlite3"
)

func newSqliteConnection(dir, dbName string, loc *time.Location) *prom.SqlConnect {
	err := os.MkdirAll(dir, 0711)
	if err != nil {
		panic(err)
	}
	return newSqlConnection("sqlite3", dir+"/"+dbName+".db", prom.FlavorSqlite, loc)
}

/*----------------------------------------------------------------------*/

const (
	sqliteTableGroup = namespace + "_group"
)

var (
	sqliteColNamesAndTypesGroup = []string{"%s VARCHAR(64)", "%s VARCHAR(255)"}
)

func sqliteInitTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sqlStm := "CREATE TABLE IF NOT EXISTS %s (" + strings.Join(sqliteColNamesAndTypesGroup, ",") + ",PRIMARY KEY (%s))"
	sqlStm = fmt.Sprintf(sqlStm, tableName, sqlColGroupId, sqlColGroupName, sqlColGroupId)
	_, err := sqlc.GetDB().Exec(sqlStm)
	if err != nil {
		panic(err)
	}
}

func newGroupDaoSqlite(sqlc *prom.SqlConnect, tableName string) GroupDao {
	return newGroupDaoSql(sqlc, tableName)
}

/*----------------------------------------------------------------------*/

const (
	sqliteTableUser = namespace + "_user"
)

var (
	sqliteColNamesAndTypesUser = []string{"%s VARCHAR(64)", "%s VARCHAR(64)", "%s VARCHAR(64)", "%s VARCHAR(64)"}
)

func sqliteInitTableUser(sqlc *prom.SqlConnect, tableName string) {
	sqlStm := "CREATE TABLE IF NOT EXISTS %s (" + strings.Join(sqliteColNamesAndTypesUser, ",") + ",PRIMARY KEY (%s))"
	sqlStm = fmt.Sprintf(sqlStm, tableName, sqlColUserUsername, sqlColUserPassword, sqlColUserName, sqlColUserGroupId, sqlColUserUsername)
	_, err := sqlc.GetDB().Exec(sqlStm)
	if err != nil {
		panic(err)
	}
}

func newUserDaoSqlite(sqlc *prom.SqlConnect, tableName string) UserDao {
	return newUserDaoSql(sqlc, tableName)
}
