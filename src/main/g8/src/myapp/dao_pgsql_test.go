package myapp

import (
	"os"
	"testing"

	"github.com/btnguyen2k/prom/sql"
)

func TestGroupDaoPgsql_GetNotExists(t *testing.T) {
	testName := "TestGroupDaoPgsql_GetNotExists"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetNotExists(t, testName, dao)
}

func TestGroupDaoPgsql_CreateGet(t *testing.T) {
	testName := "TestGroupDaoPgsql_CreateGet"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateGet(t, testName, dao)
}

func TestGroupDaoPgsql_DeleteNotExists(t *testing.T) {
	testName := "TestGroupDaoPgsql_DeleteNotExists"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoDeleteNotExists(t, testName, dao)
}

func TestGroupDaoPgsql_CreateDelete(t *testing.T) {
	testName := "TestGroupDaoPgsql_CreateDelete"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateDelete(t, testName, dao)
}

func TestGroupDaoPgsql_UpdateNotExists(t *testing.T) {
	testName := "TestGroupDaoPgsql_UpdateNotExists"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoUpdateNotExists(t, testName, dao)
}

func TestGroupDaoPgsql_CreateUpdate(t *testing.T) {
	testName := "TestGroupDaoPgsql_CreateUpdate"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateUpdate(t, testName, dao)
}

func TestGroupDaoPgsql_GetN(t *testing.T) {
	testName := "TestGroupDaoPgsql_GetN"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetN(t, testName, dao)
}

func TestGroupDaoPgsql_GetAll(t *testing.T) {
	testName := "TestGroupDaoPgsql_GetAll"
	dao := _initGroupDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetAll(t, testName, dao)
}

/*----------------------------------------------------------------------*/

func TestUserDaoPgsql_GetNotExists(t *testing.T) {
	testName := "TestUserDaoPgsql_GetNotExists"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetNotExists(t, testName, dao)
}

func TestUserDaoPgsql_CreateGet(t *testing.T) {
	testName := "TestUserDaoPgsql_CreateGet"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateGet(t, testName, dao)
}

func TestUserDaoPgsql_DeleteNotExists(t *testing.T) {
	testName := "TestUserDaoPgsql_DeleteNotExists"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoDeleteNotExists(t, testName, dao)
}

func TestUserDaoPgsql_CreateDelete(t *testing.T) {
	testName := "TestUserDaoPgsql_CreateDelete"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateDelete(t, testName, dao)
}

func TestUserDaoPgsql_UpdateNotExists(t *testing.T) {
	testName := "TestUserDaoPgsql_UpdateNotExists"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoUpdateNotExists(t, testName, dao)
}

func TestUserDaoPgsql_CreateUpdate(t *testing.T) {
	testName := "TestUserDaoPgsql_CreateUpdate"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateUpdate(t, testName, dao)
}

func TestUserDaoPgsql_GetN(t *testing.T) {
	testName := "TestUserDaoPgsql_GetN"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetN(t, testName, dao)
}

func TestUserDaoPgsql_GetAll(t *testing.T) {
	testName := "TestUserDaoPgsql_GetAll"
	dao := _initUserDao(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl), testSqlTableNameGroup, sql.FlavorPgSql)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetAll(t, testName, dao)
}
