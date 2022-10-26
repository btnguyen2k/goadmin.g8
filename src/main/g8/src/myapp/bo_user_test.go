package myapp

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/btnguyen2k/prom/sql"
)

var (
	testSqlTableNameUser        = "test_user"
	testMongoCollectionNameUser = "test_user"
)

func _initUserDaoMongo(url, db, collectionName string) UserDao {
	mc, err := _newMongoConnect(url, db)
	if err != nil {
		panic(err)
	}
	if mc == nil {
		return nil
	}
	mc.DropCollection(collectionName)
	mongoInitCollectionUser(mc, collectionName)
	return newUserDaoMongo(mc, collectionName)
}

func _initUserDaoSql(driver, url, tableName string, flavor sql.DbFlavor) UserDao {
	sqlc, err := _newSqlConnect(driver, url, testTimeZone, flavor)
	if err != nil || sqlc == nil {
		return nil
	}
	switch flavor {
	case sql.FlavorSqlite:
		sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		sqliteInitTableUser(sqlc, tableName)
		return newUserDaoSqlite(sqlc, tableName)
	case sql.FlavorMySql:
		sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		mysqlInitTableUser(sqlc, tableName)
		return newUserDaoMysql(sqlc, tableName)
	case sql.FlavorPgSql:
		sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		pgsqlInitTableUser(sqlc, tableName)
		return newUserDaoPgsql(sqlc, tableName)
	}
	sqlc.Close()
	return nil
}

func testUserDaoGetNotExists(t *testing.T, testName string, dao UserDao) {
	username := "username-notexists"
	user, err := dao.Get(username)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if user != nil {
		t.Fatalf("%s failed: expected nil", testName)
	}
}

func testUserDaoCreateGet(t *testing.T, testName string, dao UserDao) {
	username, encpwd, name, groupId := "username", encryptPassword("salt", "S3cr3t"), "User 1", "group-id"
	result, err := dao.Create(username, encpwd, name, groupId)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	user, err := dao.Get(username)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if user == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	expected := _m{"username": username, "pwd": encpwd, "name": name, "group-id": groupId}
	value := _m{"username": user.Username, "pwd": user.Password, "name": user.Name, "group-id": user.GroupId}
	if !reflect.DeepEqual(value, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, expected, value)
	}
}

func testUserDaoDeleteNotExists(t *testing.T, testName string, dao UserDao) {
	user := &User{Username: "username", Password: "encrypted-password", Name: "User 1", GroupId: "group-id"}
	result, err := dao.Delete(user)
	if result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
}

func testUserDaoCreateDelete(t *testing.T, testName string, dao UserDao) {
	username, encpwd, name, groupId := "username", encryptPassword("salt", "S3cr3t"), "User 1", "group-id"
	result, err := dao.Create(username, encpwd, name, groupId)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	user := &User{Username: username, Password: encpwd, Name: name, GroupId: groupId}
	result, err = dao.Delete(user)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	user, err = dao.Get(username)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if user != nil {
		t.Fatalf("%s failed: expected nil", testName)
	}
}

func testUserDaoUpdateNotExists(t *testing.T, testName string, dao UserDao) {
	user := &User{Username: "username-notexists", Password: "encrypted-password", Name: "User not exists", GroupId: "group-id"}
	result, err := dao.Update(user)
	if result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
}

func testUserDaoCreateUpdate(t *testing.T, testName string, dao UserDao) {
	username, encpwd, name, groupId := "username", encryptPassword("salt", "S3cr3t"), "User 1", "group-id"
	result, err := dao.Create(username, encpwd, name, groupId)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}

	newPwd := encpwd + "-new"
	newName := name + "-new"
	newGroupId := groupId + "-new"
	user := &User{Username: username, Password: newPwd, Name: newName, GroupId: newGroupId}
	result, err = dao.Update(user)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	user, err = dao.Get(username)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if user == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	expected := _m{"username": username, "pwd": newPwd, "name": newName, "group-id": newGroupId}
	value := _m{"username": user.Username, "pwd": user.Password, "name": user.Name, "group-id": user.GroupId}
	if !reflect.DeepEqual(value, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, expected, value)
	}
}

func testUserDaoGetN(t *testing.T, testName string, dao UserDao) {
	numRows := 100
	usernameList := make([]string, numRows)
	for i := 0; i < numRows; i++ {
		username, encpwd, name, groupId := fmt.Sprintf("%03d", i), encryptPassword("salt", "S3cr3t"), "User "+strconv.Itoa(i), fmt.Sprintf("group-%03d", rand.Intn(10))
		usernameList[i] = username
		result, err := dao.Create(username, encpwd, name, groupId)
		if !result || err != nil {
			t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
		}
	}

	expectedOffset := 10
	expectedNumRows := 11
	result, err := dao.GetN(expectedOffset, expectedNumRows)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(result) != expectedNumRows {
		t.Fatalf("%s failed: expected %d rows but received %d", testName, expectedNumRows, len(result))
	}
	for i, user := range result {
		if user.Username != usernameList[expectedOffset+i] {
			t.Fatalf("%s failed: expected row #%d is %s but received %s", testName, i, usernameList[expectedOffset+i], user.Username)
		}
	}
}

func testUserDaoGetAll(t *testing.T, testName string, dao UserDao) {
	numRows := 100
	usernameList := make([]string, numRows)
	for i := 0; i < numRows; i++ {
		username, encpwd, name, groupId := fmt.Sprintf("%03d", i), encryptPassword("salt", "S3cr3t"), "User "+strconv.Itoa(i), fmt.Sprintf("group-%03d", rand.Intn(10))
		usernameList[i] = username
		result, err := dao.Create(username, encpwd, name, groupId)
		if !result || err != nil {
			t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
		}
	}

	result, err := dao.GetAll()
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(result) != numRows {
		t.Fatalf("%s failed: expected %d rows but received %d", testName, numRows, len(result))
	}
	for i, user := range result {
		if user.Username != usernameList[i] {
			t.Fatalf("%s failed: expected row #%d is %s but received %s", testName, i, usernameList[i], user.Username)
		}
	}
}
