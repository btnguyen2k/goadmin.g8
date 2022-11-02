package myapp

import (
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/mongo"
	prom "github.com/btnguyen2k/prom/mongo"
)

func newMongoConnection(url, db string) *prom.MongoConnect {
	mongoConnect, err := prom.NewMongoConnect(url, db, 10000)
	if err != nil {
		panic(err)
	}
	return mongoConnect
}

/*----------------------------------------------------------------------*/

const mongoFieldId = "_id"

var (
	mongoDefaultSoringGroup = (&godal.SortingField{FieldName: mongoFieldId}).ToSortingOpt()
)

const (
	mongoCollectionGroup = namespace + "_group"
)

func mongoInitCollectionGroup(mc *prom.MongoConnect, collectionName string) {
	err := mc.CreateCollection(collectionName)
	if err != nil {
		panic(err)
	}
}

func newGroupDaoMongo(mc *prom.MongoConnect, collectionName string) GroupDao {
	dao := &GroupDaoMongo{collectionName: collectionName}
	dao.GenericDaoMongo = mongo.NewGenericDaoMongo(mc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(mongo.GenericRowMapperMongoInstance)
	if strings.Index(mc.GetUrl(), "replicaSet=") > 0 {
		dao.SetTxModeOnWrite(true)
	}
	return dao
}

/*----------------------------------------------------------------------*/

type GroupDaoMongo struct {
	*mongo.GenericDaoMongo
	collectionName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *GroupDaoMongo) GdaoCreateFilter(collectionName string, bo godal.IGenericBo) godal.FilterOpt {
	// special case for MongoDB: GBO's fieldGroupId <--> MongoDB's _id
	id, _ := bo.GboGetAttr(fieldGroupId, reddo.TypeString)
	return godal.MakeFilter(map[string]interface{}{mongoFieldId: id})
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *GroupDaoMongo) toBo(gbo godal.IGenericBo) *Group {
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
func (dao *GroupDaoMongo) toGbo(bo *Group) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(mongoFieldId, bo.Id) // special case for MongoDB
	gbo.GboSetAttr(fieldGroupId, bo.Id)
	gbo.GboSetAttr(fieldGroupName, bo.Name)
	return gbo
}

// Delete implements GroupDao.Delete
func (dao *GroupDaoMongo) Delete(bo *Group) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements GroupDao.Create
func (dao *GroupDaoMongo) Create(id, name string) (bool, error) {
	bo := &Group{
		Id:   strings.ToLower(strings.TrimSpace(id)),
		Name: strings.TrimSpace(name),
	}
	numRows, err := dao.GdaoCreate(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Get
func (dao *GroupDaoMongo) Get(id string) (*Group, error) {
	filter := godal.MakeFilter(map[string]interface{}{mongoFieldId: id})
	gbo, err := dao.GdaoFetchOne(dao.collectionName, filter)
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements GroupDao.GetN
func (dao *GroupDaoMongo) GetN(fromOffset, maxNumRows int) ([]*Group, error) {
	gboList, err := dao.GdaoFetchMany(dao.collectionName, nil, mongoDefaultSoringGroup, fromOffset, maxNumRows)
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
func (dao *GroupDaoMongo) GetAll() ([]*Group, error) {
	return dao.GetN(0, 0)
}

// Update implements GroupDao.Update
func (dao *GroupDaoMongo) Update(bo *Group) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

/*----------------------------------------------------------------------*/

const (
	mongoCollectionUser = namespace + "_user"
)

var (
	mongoDefaultSoringUser = (&godal.SortingField{FieldName: fieldUserUsername}).ToSortingOpt()
)

func mongoInitCollectionUser(mc *prom.MongoConnect, collectionName string) {
	err := mc.CreateCollection(collectionName)
	if err != nil {
		panic(err)
	}
}

func newUserDaoMongo(mc *prom.MongoConnect, collectionName string) UserDao {
	dao := &UserDaoMongo{collectionName: collectionName}
	dao.GenericDaoMongo = mongo.NewGenericDaoMongo(mc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(mongo.GenericRowMapperMongoInstance)
	if strings.Index(mc.GetUrl(), "replicaSet=") > 0 {
		dao.SetTxModeOnWrite(true)
	}
	return dao
}

type UserDaoMongo struct {
	*mongo.GenericDaoMongo
	collectionName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoMongo) GdaoCreateFilter(collectionName string, bo godal.IGenericBo) godal.FilterOpt {
	if collectionName == dao.collectionName {
		// special case for MongoDB: GBO's fieldUserUsername <--> MongoDB's _id
		username, _ := bo.GboGetAttr(fieldUserUsername, reddo.TypeString)
		return godal.MakeFilter(map[string]interface{}{mongoFieldId: username})
	}
	return nil
}

// it is recommended to have a function that transforms godal.IGenericBo to business object and vice versa.
func (dao *UserDaoMongo) toBo(gbo godal.IGenericBo) *User {
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
func (dao *UserDaoMongo) toGbo(bo *User) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(mongoFieldId, bo.Username) // special case for MongoDB
	gbo.GboSetAttr(fieldUserUsername, bo.Username)
	gbo.GboSetAttr(fieldUserPassword, bo.Password)
	gbo.GboSetAttr(fieldUserName, bo.Name)
	gbo.GboSetAttr(fieldUserGroupId, bo.GroupId)
	return gbo
}

// Delete implements UserDao.Delete
func (dao *UserDaoMongo) Delete(bo *User) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements UserDao.Create
func (dao *UserDaoMongo) Create(username, encryptedPassword, name, groupId string) (bool, error) {
	bo := &User{
		Username: strings.ToLower(strings.TrimSpace(username)),
		Password: strings.TrimSpace(encryptedPassword),
		Name:     strings.TrimSpace(name),
		GroupId:  strings.ToLower(strings.TrimSpace(groupId)),
	}
	numRows, err := dao.GdaoCreate(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements UserDao.Get
func (dao *UserDaoMongo) Get(username string) (*User, error) {
	filter := godal.MakeFilter(map[string]interface{}{mongoFieldId: username})
	gbo, err := dao.GdaoFetchOne(dao.collectionName, filter)
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoMongo) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	gboList, err := dao.GdaoFetchMany(dao.collectionName, nil, sqlDefaultSoringUser, fromOffset, maxNumRows)
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
func (dao *UserDaoMongo) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoMongo) Update(bo *User) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}
