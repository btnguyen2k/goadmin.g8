package myapp

import (
	"os"
	"testing"

	"github.com/btnguyen2k/prom/sql"
)

func TestGroupDaoSqlite_GetNotExists(t *testing.T) {
	testName := "TestGroupDaoSqlite_GetNotExists"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetNotExists(t, testName, dao)
}

func TestGroupDaoSqlite_CreateGet(t *testing.T) {
	testName := "TestGroupDaoSqlite_CreateGet"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateGet(t, testName, dao)
}

func TestGroupDaoSqlite_DeleteNotExists(t *testing.T) {
	testName := "TestGroupDaoSqlite_DeleteNotExists"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoDeleteNotExists(t, testName, dao)
}

func TestGroupDaoSqlite_CreateDelete(t *testing.T) {
	testName := "TestGroupDaoSqlite_CreateDelete"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateDelete(t, testName, dao)
}

func TestGroupDaoSqlite_UpdateNotExists(t *testing.T) {
	testName := "TestGroupDaoSqlite_UpdateNotExists"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoUpdateNotExists(t, testName, dao)
}

func TestGroupDaoSqlite_CreateUpdate(t *testing.T) {
	testName := "TestGroupDaoSqlite_CreateUpdate"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoCreateUpdate(t, testName, dao)
}

func TestGroupDaoSqlite_GetN(t *testing.T) {
	testName := "TestGroupDaoSqlite_GetN"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetN(t, testName, dao)
}

func TestGroupDaoSqlite_GetAll(t *testing.T) {
	testName := "TestGroupDaoSqlite_GetAll"
	dao := _initGroupDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoSql).GetSqlConnect().Close()
	testGroupDaoGetAll(t, testName, dao)
}

/*----------------------------------------------------------------------*/

func TestUserDaoSqlite_GetNotExists(t *testing.T) {
	testName := "TestUserDaoSqlite_GetNotExists"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetNotExists(t, testName, dao)
}

func TestUserDaoSqlite_CreateGet(t *testing.T) {
	testName := "TestUserDaoSqlite_CreateGet"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateGet(t, testName, dao)
}

func TestUserDaoSqlite_DeleteNotExists(t *testing.T) {
	testName := "TestUserDaoSqlite_DeleteNotExists"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoDeleteNotExists(t, testName, dao)
}

func TestUserDaoSqlite_CreateDelete(t *testing.T) {
	testName := "TestUserDaoSqlite_CreateDelete"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateDelete(t, testName, dao)
}

func TestUserDaoSqlite_UpdateNotExists(t *testing.T) {
	testName := "TestUserDaoSqlite_UpdateNotExists"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoUpdateNotExists(t, testName, dao)
}

func TestUserDaoSqlite_CreateUpdate(t *testing.T) {
	testName := "TestUserDaoSqlite_CreateUpdate"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoCreateUpdate(t, testName, dao)
}

func TestUserDaoSqlite_GetN(t *testing.T) {
	testName := "TestUserDaoSqlite_GetN"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetN(t, testName, dao)
}

func TestUserDaoSqlite_GetAll(t *testing.T) {
	testName := "TestUserDaoSqlite_GetAll"
	dao := _initUserDaoSql(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl), testSqlTableNameGroup, sql.FlavorSqlite)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoSql).GetSqlConnect().Close()
	testUserDaoGetAll(t, testName, dao)
}
