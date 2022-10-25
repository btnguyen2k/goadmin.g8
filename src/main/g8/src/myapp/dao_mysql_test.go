package myapp

import (
	"os"
	"testing"

	"github.com/btnguyen2k/prom/sql"
)

func TestGroupDaoMysql_GetNotExists(t *testing.T) {
	testName := "TestGroupDaoMysql_GetNotExists"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetNotExists(t, testName, dao)
}

func TestGroupDaoMysql_CreateGet(t *testing.T) {
	testName := "TestGroupDaoMysql_CreateGet"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateGet(t, testName, dao)
}

func TestGroupDaoMysql_DeleteNotExists(t *testing.T) {
	testName := "TestGroupDaoMysql_DeleteNotExists"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoDeleteNotExists(t, testName, dao)
}

func TestGroupDaoMysql_CreateDelete(t *testing.T) {
	testName := "TestGroupDaoMysql_CreateDelete"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateDelete(t, testName, dao)
}

func TestGroupDaoMysql_UpdateNotExists(t *testing.T) {
	testName := "TestGroupDaoMysql_UpdateNotExists"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoUpdateNotExists(t, testName, dao)
}

func TestGroupDaoMysql_CreateUpdate(t *testing.T) {
	testName := "TestGroupDaoMysql_CreateUpdate"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateUpdate(t, testName, dao)
}

func TestGroupDaoMysql_GetN(t *testing.T) {
	testName := "TestGroupDaoMysql_GetN"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetN(t, testName, dao)
}

func TestGroupDaoMysql_GetAll(t *testing.T) {
	testName := "TestGroupDaoMysql_GetAll"
	dao := _initGroupDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetAll(t, testName, dao)
}

/*----------------------------------------------------------------------*/

func TestUserDaoMysql_GetNotExists(t *testing.T) {
	testName := "TestUserDaoMysql_GetNotExists"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetNotExists(t, testName, dao)
}

func TestUserDaoMysql_CreateGet(t *testing.T) {
	testName := "TestUserDaoMysql_CreateGet"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateGet(t, testName, dao)
}

func TestUserDaoMysql_DeleteNotExists(t *testing.T) {
	testName := "TestUserDaoMysql_DeleteNotExists"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoDeleteNotExists(t, testName, dao)
}

func TestUserDaoMysql_CreateDelete(t *testing.T) {
	testName := "TestUserDaoMysql_CreateDelete"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateDelete(t, testName, dao)
}

func TestUserDaoMysql_UpdateNotExists(t *testing.T) {
	testName := "TestUserDaoMysql_UpdateNotExists"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoUpdateNotExists(t, testName, dao)
}

func TestUserDaoMysql_CreateUpdate(t *testing.T) {
	testName := "TestUserDaoMysql_CreateUpdate"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateUpdate(t, testName, dao)
}

func TestUserDaoMysql_GetN(t *testing.T) {
	testName := "TestUserDaoMysql_GetN"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetN(t, testName, dao)
}

func TestUserDaoMysql_GetAll(t *testing.T) {
	testName := "TestUserDaoMysql_GetAll"
	dao := _initUserDaoSql(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl), testSqlTableNameGroup, sql.FlavorMySql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetAll(t, testName, dao)
}
