package myapp

// @availabble since template-v0.4.r2

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

func newPgsqlConnection(url, timezone string) *prom.SqlConnect {
	driver := "postgres"
	sqlConnect, err := prom.NewSqlConnect(driver, url, 10000, nil)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation(timezone)
	sqlConnect.SetLocation(loc)
	sqlConnect.SetDbFlavor(prom.FlavorPgSql)
	return sqlConnect
}

func pgsqlInitTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(255), PRIMARY KEY (%s))",
		tableName, pgsqlColGroupId, pgsqlColGroupName, pgsqlColGroupId)
	_, err := sqlc.GetDB().Exec(sql)
	if err != nil {
		panic(err)
	}
}

func pgsqlInitTableUser(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(64), %s VARCHAR(64), %s VARCHAR(64), PRIMARY KEY (%s))",
		tableName, pgsqlColUserUsername, pgsqlColUserPassword, pgsqlColUserName, pgsqlColUserGroupId, pgsqlColUserUsername)
	_, err := sqlc.GetDB().Exec(sql)
	if err != nil {
		panic(err)
	}
}

/*----------------------------------------------------------------------*/

func newUserDaoPgsql(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoPgsql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: pgsqlMapFieldToColNameUser},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: pgsqlMapColNameToFieldUser},
		ColumnsListMap:              map[string][]string{tableName: pgsqlColsUser},
	})
	dao.SetSqlFlavor(prom.FlavorPgSql)
	return dao
}

const (
	pgsqlTableUser       = namespace + "_user"
	pgsqlColUserUsername = "uname"
	pgsqlColUserPassword = "upwd"
	pgsqlColUserName     = "display_name"
	pgsqlColUserGroupId  = "gid"
)

var (
	pgsqlColsUser              = []string{pgsqlColUserUsername, pgsqlColUserPassword, pgsqlColUserName, pgsqlColUserGroupId}
	pgsqlMapFieldToColNameUser = map[string]interface{}{fieldUserUsername: pgsqlColUserUsername, fieldUserPassword: pgsqlColUserPassword, fieldUserName: pgsqlColUserName, fieldUserGroupId: pgsqlColUserGroupId}
	pgsqlMapColNameToFieldUser = map[string]interface{}{pgsqlColUserUsername: fieldUserUsername, pgsqlColUserPassword: fieldUserPassword, pgsqlColUserName: fieldUserName, pgsqlColUserGroupId: fieldUserGroupId}
)

type UserDaoPgsql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoPgsql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{pgsqlColUserUsername: bo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString)}
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *UserDaoPgsql) toBo(gbo godal.IGenericBo) *User {
	if gbo == nil {
		return nil
	}
	bo := &User{
		Username: gbo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString).(string),
		Password: gbo.GboGetAttrUnsafe(fieldUserPassword, reddo.TypeString).(string),
		Name:     gbo.GboGetAttrUnsafe(fieldUserName, reddo.TypeString).(string),
		GroupId:  gbo.GboGetAttrUnsafe(fieldUserGroupId, reddo.TypeString).(string),
	}
	return bo
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *UserDaoPgsql) toGbo(bo *User) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldUserUsername, bo.Username)
	gbo.GboSetAttr(fieldUserPassword, bo.Password)
	gbo.GboSetAttr(fieldUserName, bo.Name)
	gbo.GboSetAttr(fieldUserGroupId, bo.GroupId)
	return gbo
}

// Delete implements UserDao.Delete
func (dao *UserDaoPgsql) Delete(bo *User) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements UserDao.Create
func (dao *UserDaoPgsql) Create(username, encryptedPassword, name, groupId string) (bool, error) {
	bo := &User{
		Username: strings.ToLower(strings.TrimSpace(username)),
		Password: strings.TrimSpace(encryptedPassword),
		Name:     strings.TrimSpace(name),
		GroupId:  strings.ToLower(strings.TrimSpace(groupId)),
	}
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements UserDao.Get
func (dao *UserDaoPgsql) Get(username string) (*User, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{pgsqlColUserUsername: username})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoPgsql) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, nil, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*User, 0)
	for _, gbo := range gboList {
		bo := dao.toBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements UserDao.GetAll
func (dao *UserDaoPgsql) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoPgsql) Update(bo *User) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

/*----------------------------------------------------------------------*/

func newGroupDaoPgsql(sqlc *prom.SqlConnect, tableName string) GroupDao {
	dao := &GroupDaoPgsql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: pgsqlMapFieldToColNameGroup},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: pgsqlMapColNameToFieldGroup},
		ColumnsListMap:              map[string][]string{tableName: pgsqlColsGroup},
	})
	dao.SetSqlFlavor(prom.FlavorPgSql)
	return dao
}

const (
	pgsqlTableGroup   = namespace + "_group"
	pgsqlColGroupId   = "gid"
	pgsqlColGroupName = "gname"
)

var (
	pgsqlColsGroup              = []string{pgsqlColGroupId, pgsqlColGroupName}
	pgsqlMapFieldToColNameGroup = map[string]interface{}{fieldGroupId: pgsqlColGroupId, fieldGroupName: pgsqlColGroupName}
	pgsqlMapColNameToFieldGroup = map[string]interface{}{pgsqlColGroupId: fieldGroupId, pgsqlColGroupName: fieldGroupName}
)

type GroupDaoPgsql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *GroupDaoPgsql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{pgsqlColGroupId: bo.GboGetAttrUnsafe(fieldGroupId, reddo.TypeString)}
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *GroupDaoPgsql) toBo(gbo godal.IGenericBo) *Group {
	if gbo == nil {
		return nil
	}
	bo := &Group{
		Id:   gbo.GboGetAttrUnsafe(fieldGroupId, reddo.TypeString).(string),
		Name: gbo.GboGetAttrUnsafe(fieldGroupName, reddo.TypeString).(string),
	}
	return bo
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *GroupDaoPgsql) toGbo(bo *Group) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldGroupId, bo.Id)
	gbo.GboSetAttr(fieldGroupName, bo.Name)
	return gbo
}

// Delete implements GroupDao.Delete
func (dao *GroupDaoPgsql) Delete(bo *Group) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Create
func (dao *GroupDaoPgsql) Create(id, name string) (bool, error) {
	bo := &Group{
		Id:   strings.ToLower(strings.TrimSpace(id)),
		Name: strings.TrimSpace(name),
	}
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Get
func (dao *GroupDaoPgsql) Get(id string) (*Group, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{pgsqlColGroupId: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements GroupDao.GetN
func (dao *GroupDaoPgsql) GetN(fromOffset, maxNumRows int) ([]*Group, error) {
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, nil, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*Group, 0)
	for _, gbo := range gboList {
		bo := dao.toBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements GroupDao.GetAll
func (dao *GroupDaoPgsql) GetAll() ([]*Group, error) {
	return dao.GetN(0, 0)
}

// Update implements GroupDao.Update
func (dao *GroupDaoPgsql) Update(bo *Group) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
