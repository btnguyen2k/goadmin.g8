package myapp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/btnguyen2k/prom/sql"
)

var (
	testSqlTableNameGroup = "test_group"
)

func _initGroupDao(driver, url, tableName string, flavor sql.DbFlavor) GroupDao {
	sqlc, err := _newSqlConnect(driver, url, testTimeZone, flavor)
	if err != nil {
		panic(err)
	}
	if sqlc == nil {
		return nil
	}
	switch flavor {
	case sql.FlavorSqlite:
		sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		sqliteInitTableGroup(sqlc, tableName)
		return newGroupDaoSqlite(sqlc, tableName)
	case sql.FlavorPgSql:
		sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		pgsqlInitTableGroup(sqlc, tableName)
		return newGroupDaoPgsql(sqlc, tableName)
	}
	sqlc.Close()
	return nil
}

func testGroupDaoGetNotExists(t *testing.T, testName string, dao GroupDao) {
	groupId := "group-id-notexists"
	group, err := dao.Get(groupId)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if group != nil {
		t.Fatalf("%s failed: expected nil", testName)
	}
}

func testGroupDaoCreateGet(t *testing.T, testName string, dao GroupDao) {
	groupId, groupName := "group-id", "group-name"
	result, err := dao.Create(groupId, groupName)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	group, err := dao.Get(groupId)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if group == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	if group.Id != groupId || group.Name != groupName {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, _m{"id": groupId, "name": groupName}, _m{"id": group.Id, "name": group.Name})
	}
}

func testGroupDaoDeleteNotExists(t *testing.T, testName string, dao GroupDao) {
	group := &Group{Id: "group-id-notexists", Name: "name"}
	result, err := dao.Delete(group)
	if result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
}

func testGroupDaoCreateDelete(t *testing.T, testName string, dao GroupDao) {
	groupId, groupName := "group-id", "group-name"
	result, err := dao.Create(groupId, groupName)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	group := &Group{Id: groupId, Name: groupName}
	result, err = dao.Delete(group)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	group, err = dao.Get(groupId)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if group != nil {
		t.Fatalf("%s failed: expected nil", testName)
	}
}

func testGroupDaoUpdateNotExists(t *testing.T, testName string, dao GroupDao) {
	group := &Group{Id: "group-id-notexists", Name: "name"}
	result, err := dao.Update(group)
	if result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
}

func testGroupDaoCreateUpdate(t *testing.T, testName string, dao GroupDao) {
	groupId, groupName := "group-id", "group-name"
	result, err := dao.Create(groupId, groupName)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}

	newGroupName := groupName + "-new"
	group := &Group{Id: groupId, Name: newGroupName}
	result, err = dao.Update(group)
	if !result || err != nil {
		t.Fatalf("%s failed: {result %#v / error %s}", testName, result, err)
	}
	group, err = dao.Get(groupId)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if group == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	if group.Id != groupId || group.Name != newGroupName {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, _m{"id": groupId, "name": newGroupName}, _m{"id": group.Id, "name": group.Name})
	}
}

func testGroupDaoGetN(t *testing.T, testName string, dao GroupDao) {
	numRows := 100
	idList := make([]string, 100)
	for i := 0; i < numRows; i++ {
		groupId := fmt.Sprintf("%03d", i)
		idList[i] = groupId
		groupName := "group-name-" + strconv.Itoa(i)
		result, err := dao.Create(groupId, groupName)
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
	for i, group := range result {
		if group.Id != idList[expectedOffset+i] {
			t.Fatalf("%s failed: expected row #%d is %s but received %s", testName, i, idList[expectedOffset+i], group.Id)
		}
	}
}

func testGroupDaoGetAll(t *testing.T, testName string, dao GroupDao) {
	numRows := 100
	idList := make([]string, 100)
	for i := 0; i < numRows; i++ {
		groupId := fmt.Sprintf("%03d", i)
		idList[i] = groupId
		groupName := "group-name-" + strconv.Itoa(i)
		result, err := dao.Create(groupId, groupName)
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
	for i, group := range result {
		if group.Id != idList[i] {
			t.Fatalf("%s failed: expected row #%d is %s but received %s", testName, i, idList[i], group.Id)
		}
	}
}
