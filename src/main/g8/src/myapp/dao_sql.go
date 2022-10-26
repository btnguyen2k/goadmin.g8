package myapp

import (
	"strings"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	prom "github.com/btnguyen2k/prom/sql"
)

func newSqlConnection(driver, dsn string, flavor prom.DbFlavor, loc *time.Location) *prom.SqlConnect {
	sqlConnect, err := prom.NewSqlConnectWithFlavor(driver, dsn, 10000, nil, flavor)
	if err != nil {
		panic(err)
	}
	if loc != nil {
		sqlConnect.SetLocation(loc)
	} else {
		sqlConnect.SetLocation(time.UTC)
	}
	return sqlConnect
}

/*----------------------------------------------------------------------*/

const (
	sqlColGroupId   = "gid"
	sqlColGroupName = "gname"
)

var (
	sqlColsGroup              = []string{sqlColGroupId, sqlColGroupName}
	sqlMapFieldToColNameGroup = map[string]interface{}{fieldGroupId: sqlColGroupId, fieldGroupName: sqlColGroupName}
	sqlMapColNameToFieldGroup = map[string]interface{}{sqlColGroupId: fieldGroupId, sqlColGroupName: fieldGroupName}
	sqlDefaultSoringGroup     = (&godal.SortingOpt{}).Add(&godal.SortingField{FieldName: fieldGroupId})
)

func newGroupDaoSql(sqlc *prom.SqlConnect, tableName string) GroupDao {
	dao := &GroupDaoSql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: sqlMapFieldToColNameGroup},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: sqlMapColNameToFieldGroup},
		ColumnsListMap:              map[string][]string{tableName: sqlColsGroup},
	})
	return dao
}

type GroupDaoSql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *GroupDaoSql) GdaoCreateFilter(tableName string, bo godal.IGenericBo) godal.FilterOpt {
	id, _ := bo.GboGetAttr(fieldGroupId, reddo.TypeString)
	return &godal.FilterOptFieldOpValue{FieldName: fieldGroupId, Operator: godal.FilterOpEqual, Value: id}
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *GroupDaoSql) toBo(gbo godal.IGenericBo) *Group {
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
func (dao *GroupDaoSql) toGbo(bo *Group) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldGroupId, bo.Id)
	gbo.GboSetAttr(fieldGroupName, bo.Name)
	return gbo
}

// Delete implements GroupDao.Delete
func (dao *GroupDaoSql) Delete(bo *Group) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements GroupDao.Create
func (dao *GroupDaoSql) Create(id, name string) (bool, error) {
	bo := &Group{
		Id:   strings.ToLower(strings.TrimSpace(id)),
		Name: strings.TrimSpace(name),
	}
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Get
func (dao *GroupDaoSql) Get(id string) (*Group, error) {
	filter := &godal.FilterOptFieldOpValue{FieldName: fieldGroupId, Operator: godal.FilterOpEqual, Value: id}
	gbo, err := dao.GdaoFetchOne(dao.tableName, filter)
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements GroupDao.GetN
func (dao *GroupDaoSql) GetN(fromOffset, maxNumRows int) ([]*Group, error) {
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, sqlDefaultSoringGroup, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*Group, len(gboList))
	for i, gbo := range gboList {
		result[i] = dao.toBo(gbo)
	}
	return result, nil
}

// GetAll implements GroupDao.GetAll
func (dao *GroupDaoSql) GetAll() ([]*Group, error) {
	return dao.GetN(0, 0)
}

// Update implements GroupDao.Update
func (dao *GroupDaoSql) Update(bo *Group) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

/*----------------------------------------------------------------------*/

const (
	sqlColUserUsername = "uname"
	sqlColUserPassword = "upwd"
	sqlColUserName     = "display_name"
	sqlColUserGroupId  = "gid"
)

var (
	sqlColsUser              = []string{sqlColUserUsername, sqlColUserPassword, sqlColUserName, sqlColUserGroupId}
	sqlMapFieldToColNameUser = map[string]interface{}{fieldUserUsername: sqlColUserUsername, fieldUserPassword: sqlColUserPassword, fieldUserName: sqlColUserName, fieldUserGroupId: sqlColUserGroupId}
	sqlMapColNameToFieldUser = map[string]interface{}{sqlColUserUsername: fieldUserUsername, sqlColUserPassword: fieldUserPassword, sqlColUserName: fieldUserName, sqlColUserGroupId: fieldUserGroupId}
	sqlDefaultSoringUser     = (&godal.SortingOpt{}).Add(&godal.SortingField{FieldName: fieldUserUsername})
)

func newUserDaoSql(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoSql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: sqlMapFieldToColNameUser},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: sqlMapColNameToFieldUser},
		ColumnsListMap:              map[string][]string{tableName: sqlColsUser},
	})
	return dao
}

type UserDaoSql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoSql) GdaoCreateFilter(tableName string, bo godal.IGenericBo) godal.FilterOpt {
	if tableName == dao.tableName {
		username, _ := bo.GboGetAttr(fieldUserUsername, reddo.TypeString)
		return &godal.FilterOptFieldOpValue{FieldName: fieldUserUsername, Operator: godal.FilterOpEqual, Value: username}
	}
	return nil
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *UserDaoSql) toBo(gbo godal.IGenericBo) *User {
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
func (dao *UserDaoSql) toGbo(bo *User) godal.IGenericBo {
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
func (dao *UserDaoSql) Delete(bo *User) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements UserDao.Create
func (dao *UserDaoSql) Create(username, encryptedPassword, name, groupId string) (bool, error) {
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
func (dao *UserDaoSql) Get(username string) (*User, error) {
	filter := &godal.FilterOptFieldOpValue{FieldName: fieldUserUsername, Operator: godal.FilterOpEqual, Value: username}
	gbo, err := dao.GdaoFetchOne(dao.tableName, filter)
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoSql) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, sqlDefaultSoringUser, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*User, len(gboList))
	for i, gbo := range gboList {
		result[i] = dao.toBo(gbo)
	}
	return result, nil
}

// GetAll implements UserDao.GetAll
func (dao *UserDaoSql) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoSql) Update(bo *User) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
