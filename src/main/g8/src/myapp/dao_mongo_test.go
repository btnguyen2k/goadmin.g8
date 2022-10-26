package myapp

import (
	"os"
	"testing"
)

func TestGroupDaoMongo_GetNotExists(t *testing.T) {
	testName := "TestGroupDaoMongo_GetNotExists"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoGetNotExists(t, testName, dao)
}

func TestGroupDaoMongo_CreateGet(t *testing.T) {
	testName := "TestGroupDaoMongo_CreateGet"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoCreateGet(t, testName, dao)
}

func TestGroupDaoMongo_DeleteNotExists(t *testing.T) {
	testName := "TestGroupDaoMongo_DeleteNotExists"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoDeleteNotExists(t, testName, dao)
}

func TestGroupDaoMongo_CreateDelete(t *testing.T) {
	testName := "TestGroupDaoMongo_CreateDelete"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoCreateDelete(t, testName, dao)
}

func TestGroupDaoMongo_UpdateNotExists(t *testing.T) {
	testName := "TestGroupDaoMongo_UpdateNotExists"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoUpdateNotExists(t, testName, dao)
}

func TestGroupDaoMongo_CreateUpdate(t *testing.T) {
	testName := "TestGroupDaoMongo_CreateUpdate"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoCreateUpdate(t, testName, dao)
}

func TestGroupDaoMongo_GetN(t *testing.T) {
	testName := "TestGroupDaoMongo_GetN"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoGetN(t, testName, dao)
}

func TestGroupDaoMongo_GetAll(t *testing.T) {
	testName := "TestGroupDaoMongo_GetAll"
	dao := _initGroupDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameGroup)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*GroupDaoMongo).GetMongoConnect().Close(nil)
	testGroupDaoGetAll(t, testName, dao)
}

/*----------------------------------------------------------------------*/

func TestUserDaoMongo_GetNotExists(t *testing.T) {
	testName := "TestUserDaoMongo_GetNotExists"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoGetNotExists(t, testName, dao)
}

func TestUserDaoMongo_CreateGet(t *testing.T) {
	testName := "TestUserDaoMongo_CreateGet"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoCreateGet(t, testName, dao)
}

func TestUserDaoMongo_DeleteNotExists(t *testing.T) {
	testName := "TestUserDaoMongo_DeleteNotExists"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoDeleteNotExists(t, testName, dao)
}

func TestUserDaoMongo_CreateDelete(t *testing.T) {
	testName := "TestUserDaoMongo_CreateDelete"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoCreateDelete(t, testName, dao)
}

func TestUserDaoMongo_UpdateNotExists(t *testing.T) {
	testName := "TestUserDaoMongo_UpdateNotExists"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoUpdateNotExists(t, testName, dao)
}

func TestUserDaoMongo_CreateUpdate(t *testing.T) {
	testName := "TestUserDaoMongo_CreateUpdate"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoCreateUpdate(t, testName, dao)
}

func TestUserDaoMongo_GetN(t *testing.T) {
	testName := "TestUserDaoMongo_GetN"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoGetN(t, testName, dao)
}

func TestUserDaoMongo_GetAll(t *testing.T) {
	testName := "TestUserDaoMongo_GetAll"
	dao := _initUserDaoMongo(os.Getenv(envMongoUrl), os.Getenv(envMongoDb), testMongoCollectionNameUser)
	if dao == nil {
		t.SkipNow()
	}
	defer dao.(*UserDaoMongo).GetMongoConnect().Close(nil)
	testUserDaoGetAll(t, testName, dao)
}
