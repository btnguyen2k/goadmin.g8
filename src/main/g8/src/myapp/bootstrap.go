// Package myapp contains application's source code.
package myapp

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/goyai"
	prommongo "github.com/btnguyen2k/prom/mongo"
	promsql "github.com/btnguyen2k/prom/sql"
	"github.com/go-akka/configuration"
	"github.com/labstack/echo/v4"
	"main/src/goadmin"
	"main/src/utils"
)

type MyBootstrapper struct {
	name string
}

var (
	Bootstrapper = &MyBootstrapper{name: "myapp"}

	demoMode     = false
	cdnMode      = false
	myStaticPath = "/static"
	sqlc         *promsql.SqlConnect
	mc           *prommongo.MongoConnect
	groupDao     GroupDao
	userDao      UserDao
	myI18n       goyai.I18n
)

const (
	namespace = "myapp"

	ctxCurrentUser = "usr"
	ctxLocale      = "loc"
	cookieLocale   = "loc"
	sessionMyUid   = "uid"

	actionNameHome          = "home"
	actionNameCpLogin       = "cp_login"
	actionNameCpLoginSubmit = "cp_login_submit"
	actionNameCpLogout      = "cp_logout"
	actionNameCpDashboard   = "cp_dashboard"
	actionNameCpProfile     = "cp_profile"

	actionNameCpChangePassword       = "cp_change_password"
	actionNameCpChangePasswordSubmit = "cp_change_password_submit"

	actionNameCpGroups            = "cp_groups"
	actionNameCpCreateGroup       = "cp_create_group"
	actionNameCpCreateGroupSubmit = "cp_create_group_submit"
	actionNameCpEditGroup         = "cp_edit_group"
	actionNameCpEditGroupSubmit   = "cp_edit_group_submit"
	actionNameCpDeleteGroup       = "cp_delete_group"
	actionNameCpDeleteGroupSubmit = "cp_delete_group_submit"

	actionNameCpUsers            = "cp_users"
	actionNameCpCreateUser       = "cp_create_user"
	actionNameCpCreateUserSubmit = "cp_create_user_submit"
	actionNameCpEditUser         = "cp_edit_user"
	actionNameCpEditUserSubmit   = "cp_edit_user_submit"
	actionNameCpDeleteUser       = "cp_delete_user"
	actionNameCpDeleteUserSubmit = "cp_delete_user_submit"
)

// Bootstrap implements goadmin.IBootstrapper.Bootstrap
//
// Bootstrapper usually does:
// - register URI routing
// - other initializing work (e.g. creating DAO, initializing database, etc)
func (b *MyBootstrapper) Bootstrap(conf *configuration.Config, e *echo.Echo) error {
	cdnMode = conf.GetBoolean(namespace+".cdn_mode", cdnMode)
	demoMode = conf.GetBoolean(namespace+".demo_mode", demoMode)
	systemUserUsername = conf.GetString(namespace+".init.admin_username", systemUserUsername)
	systemUserName = conf.GetString(namespace+".init.admin_name", systemUserName)

	myStaticPath = "/static_v" + conf.GetString("app.version", "")
	e.Static(myStaticPath, "public")

	if i18n, err := goyai.BuildI18n(goyai.I18nOptions{
		ConfigFileOrDir: "./config/i18n_" + namespace,
		DefaultLocale:   "en",
		I18nFileFormat:  goyai.Auto,
	}); err != nil {
		return err
	} else {
		myI18n = i18n
	}

	initDaos()
	_initData()

	// register a custom namespace-scope template renderer
	goadmin.EchoRegisterRenderer(namespace, newTemplateRenderer("./views/myapp", ".html"))

	e.Use(middlewarePopulateLocale)

	e.GET("/", actionHome).Name = actionNameHome

	e.GET("/cp/login", actionCpLogin).Name = actionNameCpLogin
	e.POST("/cp/login", actionCpLoginSubmit).Name = actionNameCpLoginSubmit
	e.GET("/cp/logout", actionCpLogout).Name = actionNameCpLogout
	e.GET("/cp", actionCpDashboard, middlewareRequiredAuth).Name = actionNameCpDashboard
	e.GET("/cp/profile", actionCpProfile, middlewareRequiredAuth).Name = actionNameCpProfile
	e.GET("/cp/changePassword", actionCpChangePassword, middlewareRequiredAuth).Name = actionNameCpChangePassword
	e.POST("/cp/changePassword", actionCpChangePasswordSubmit, middlewareRequiredAuth).Name = actionNameCpChangePasswordSubmit

	e.GET("/cp/groups", actionCpGroupList, middlewareRequiredAuth).Name = actionNameCpGroups
	e.GET("/cp/createGroup", actionCpCreateGroup, middlewareRequiredAuth).Name = actionNameCpCreateGroup
	e.POST("/cp/createGroup", actionCpCreateGroupSubmit, middlewareRequiredAuth).Name = actionNameCpCreateGroupSubmit
	e.GET("/cp/editGroup", actionCpEditGroup, middlewareRequiredAuth).Name = actionNameCpEditGroup
	e.POST("/cp/editGroup", actionCpEditGroupSubmit, middlewareRequiredAuth).Name = actionNameCpEditGroupSubmit
	e.GET("/cp/deleteGroup", actionCpDeleteGroup, middlewareRequiredAuth).Name = actionNameCpDeleteGroup
	e.POST("/cp/deleteGroup", actionCpDeleteGroupSubmit, middlewareRequiredAuth).Name = actionNameCpDeleteGroupSubmit

	e.GET("/cp/users", actionCpUserList, middlewareRequiredAuth).Name = actionNameCpUsers
	e.GET("/cp/createUser", actionCpCreateUser, middlewareRequiredAuth).Name = actionNameCpCreateUser
	e.POST("/cp/createUser", actionCpCreateUserSubmit, middlewareRequiredAuth).Name = actionNameCpCreateUserSubmit
	e.GET("/cp/editUser", actionCpEditUser, middlewareRequiredAuth).Name = actionNameCpEditUser
	e.POST("/cp/editUser", actionCpEditUserSubmit, middlewareRequiredAuth).Name = actionNameCpEditUserSubmit
	e.GET("/cp/deleteUser", actionCpDeleteUser, middlewareRequiredAuth).Name = actionNameCpDeleteUser
	e.POST("/cp/deleteUser", actionCpDeleteUserSubmit, middlewareRequiredAuth).Name = actionNameCpDeleteUserSubmit

	return nil
}

func initDaos() {
	dbtype := goadmin.AppConfig.GetString(namespace + ".db.type")
	switch dbtype {
	case "mongo", "mongodb":
		url := goadmin.AppConfig.GetString(namespace+".db.mongodb.url", "mongodb://test:test@localhost:37017/?authSource=admin")
		db := goadmin.AppConfig.GetString(namespace+".db.mongodb.db", "test")
		mc = newMongoConnection(url, db)
		mongoInitCollectionGroup(mc, mongoCollectionGroup)
		mongoInitCollectionUser(mc, mongoCollectionUser)
		groupDao = newGroupDaoMongo(mc, mongoCollectionGroup)
		userDao = newUserDaoMongo(mc, mongoCollectionUser)
	case "mysql":
		url := goadmin.AppConfig.GetString(namespace+".db.mysql.url", "test:test@tcp(localhost:3306)/test?charset=utf8mb4,utf8&parseTime=true&loc=${loc}")
		urlTimezone := strings.ReplaceAll(utils.Location.String(), "/", "%2f")
		url = strings.ReplaceAll(url, "${loc}", urlTimezone)
		url = strings.ReplaceAll(url, "${tz}", urlTimezone)
		url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
		sqlc = newMysqlConnection(url, utils.Location)
		mysqlInitTableGroup(sqlc, mysqlTableGroup)
		mysqlInitTableUser(sqlc, mysqlTableUser)
		groupDao = newGroupDaoMysql(sqlc, mysqlTableGroup)
		userDao = newUserDaoMysql(sqlc, mysqlTableUser)
	case "postgresql", "pgsql", "postgres":
		url := goadmin.AppConfig.GetString(namespace+".db.pgsql.url", "postgres://test:test@localhost:5432/test")
		sqlc = newPgsqlConnection(url, utils.Location)
		pgsqlInitTableGroup(sqlc, pgsqlTableGroup)
		pgsqlInitTableUser(sqlc, pgsqlTableUser)
		groupDao = newGroupDaoPgsql(sqlc, pgsqlTableGroup)
		userDao = newUserDaoPgsql(sqlc, pgsqlTableUser)
	case "sqlite", "sqlite3":
		root := goadmin.AppConfig.GetString(namespace+".db.sqlite.root", "./data/sqlite")
		sqlc = newSqliteConnection(root, namespace, utils.Location)
		sqliteInitTableGroup(sqlc, sqliteTableGroup)
		sqliteInitTableUser(sqlc, sqliteTableUser)
		groupDao = newGroupDaoSqlite(sqlc, sqliteTableGroup)
		userDao = newUserDaoSqlite(sqlc, sqliteTableUser)
	default:
		panic(fmt.Sprintf("unsupported database type: %s", dbtype))
	}
}

func _initData() {
	if systemGroup, err := groupDao.Get(systemGroupId); err != nil {
		panic("error while getting group [" + systemGroupId + "]: " + err.Error())
	} else if systemGroup == nil {
		log.Printf("System group [%s] not found, creating one...", systemGroupId)
		result, err := groupDao.Create(systemGroupId, "System User Group")
		if err != nil {
			panic("error while creating group [" + systemGroupId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create group [%s]", systemGroupId)
		}
	}

	adminUser, err := userDao.Get(systemUserUsername)
	if err != nil {
		panic("error while getting user [" + systemUserUsername + "]: " + err.Error())
	}
	if adminUser == nil {
		pwd := goadmin.AppConfig.GetString(namespace+".init.admin_password", "S3cr3t")
		log.Printf("Admin user [%s] not found, creating one with password [%s]...", systemUserName, pwd)
		result, err := userDao.Create(systemUserUsername, encryptPassword(systemUserUsername, pwd), systemUserName, systemGroupId)
		if err != nil {
			panic("error while creating user [" + systemUserUsername + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create user [%s]", systemUserUsername)
		}
	}
}

/*----------------------------------------------------------------------*/
func newTemplateRenderer(directory, templateFileSuffix string) *myRenderer {
	return &myRenderer{
		directory:          directory,
		templateFileSuffix: templateFileSuffix,
		templates:          map[string]*template.Template{},
	}
}

// myRenderer is a custom html/template renderer for Echo framework
// See: https://echo.labstack.com/guide/templates
type myRenderer struct {
	directory          string
	templateFileSuffix string
	templates          map[string]*template.Template
}

// Render renders a template document.
// - tplNames is list of template names, separated by colon (e.g. <template-name-1>[:<template-name-2>[:<template-name-3>...]])
func (r *myRenderer) Render(w io.Writer, tplNames string, data interface{}, c echo.Context) error {
	if utils.DevMode {
		log.Printf("[DEBUG] %s renderer: rendering [%s]...", namespace, tplNames)
	}

	v := reflect.ValueOf(data)
	if data == nil || v.IsNil() {
		data = make(map[string]interface{})
	}

	sess := getSession(c)
	flash := sess.Flashes()
	sess.Save(c.Request(), c.Response())

	// add global data/methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["cdn_mode"] = cdnMode
		viewContext["static"] = myStaticPath
		viewContext["i18n"] = myI18n
		viewContext["locale"] = getContextString(c, ctxLocale)
		viewContext["reverse"] = c.Echo().Reverse
		viewContext["appInfo"] = goadmin.AppConfig.GetConfig("app")
		viewContext["appUtils"] = &MyAppUtils{c: c}
		if len(flash) > 0 {
			flashMsg := flash[0].(string)
			if strings.HasPrefix(flashMsg, flashPrefixWarning) {
				viewContext["flashWarning"] = flashMsg[len(flashPrefixWarning):]
			} else if strings.HasPrefix(flashMsg, flashPrefixError) {
				viewContext["flashError"] = flashMsg[len(flashPrefixError):]
			} else if strings.HasPrefix(flashMsg, flashPrefixInfo) {
				viewContext["flashInfo"] = flashMsg[len(flashPrefixInfo):]
			} else {
				viewContext["flashInfo"] = flashMsg
			}
		}
		u := c.Get(ctxCurrentUser)
		if u != nil {
			switch u.(type) {
			case User:
				usr := u.(User)
				viewContext["currentUser"] = toUserModel(c, &usr)
			case *User:
				viewContext["currentUser"] = toUserModel(c, u.(*User))
			}
		}
	}

	tpl := r.templates[tplNames]
	tokens := strings.Split(tplNames, ":")
	if tpl == nil {
		var files []string
		for _, v := range tokens {
			files = append(files, r.directory+"/"+v+r.templateFileSuffix)
		}
		tpl = template.Must(template.New(tplNames).ParseFiles(files...))
		if !utils.DevMode {
			// DEV mode: disable template caching
			r.templates[tplNames] = tpl
		}
	}
	// first template-tplNames should be "master" template, and its tplNames is prefixed with ".html"
	return tpl.ExecuteTemplate(w, tokens[0]+".html", data)
}

/*----------------------------------------------------------------------*/
// middleware function that populate the value of "locale" field to echo.Context
// available since template-r3
func middlewarePopulateLocale(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(ctxLocale, getCookieString(c, cookieLocale))
		return next(c)
	}
}

// authentication middleware
func middlewareRequiredAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := getSession(c)
		var currentUser *User = nil
		var err error
		if uid, has := sess.Values[sessionMyUid]; has {
			uid, _ = reddo.ToString(uid)
			if uid != nil {
				username := uid.(string)
				currentUser, err = userDao.Get(username)
				if err != nil {
					log.Printf("error while fetching user [%s]: %s", username, err.Error())
				}
			}
		}
		if currentUser == nil {
			return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpLogin))
		}
		c.Set(ctxCurrentUser, currentUser)
		return next(c)
	}
}

func actionHome(c echo.Context) error {
	return c.Render(http.StatusOK, namespace+":landing", nil)
}

func actionCpLogin(c echo.Context) error {
	data := map[string]interface{}{}
	if utils.DevMode {
		formData := url.Values{
			"username": []string{systemUserUsername},
			"password": []string{goadmin.AppConfig.GetString(namespace + ".init.admin_password")},
		}
		data["form"] = formData
	}
	return c.Render(http.StatusOK, namespace+":login", data)
}

func actionCpLoginSubmit(c echo.Context) error {
	const (
		formFieldUsername = "username"
		formFieldPassword = "password"
	)
	var username, password, encPassword string
	var user *User
	var errMsg string
	var err error
	formData, err := c.FormParams()
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_form_400", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": err.Error()},
		})
		goto end
	}
	username = formData.Get(formFieldUsername)
	user, err = userDao.Get(username)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": username + "/" + err.Error()},
		})
		goto end
	}
	if user == nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_user_not_found", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"user": username},
		})
		goto end
	}
	password = formData.Get(formFieldPassword)
	encPassword = encryptPassword(user.Username, password)
	if encPassword != user.Password {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_signin_failed")
		goto end
	}

	// login successful
	setSessionValue(c, sessionMyUid, user.Username)
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpDashboard))
end:
	if utils.DevMode {
		formData.Set("username", systemUserUsername)
		formData.Set("password", goadmin.AppConfig.GetString(namespace+".init.admin_password"))
	}
	return c.Render(http.StatusOK, namespace+":login", map[string]interface{}{
		"form":  formData,
		"error": errMsg,
	})
}

func actionCpLogout(c echo.Context) error {
	setSessionValue(c, sessionMyUid, nil)
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpDashboard))
}

func actionCpDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, namespace+":layout:cp_dashboard", map[string]interface{}{
		"active":  "dashboard",
		"osUtils": &OsUtils{},
	})
}

func actionCpProfile(c echo.Context) error {
	return c.Render(http.StatusOK, namespace+":layout:cp_profile", map[string]interface{}{
		"active": "profile",
	})
}

func actionCpChangePassword(c echo.Context) error {
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpProfile))
}

func actionCpChangePasswordSubmit(c echo.Context) error {
	var encPwd, currentPwd, pwd, pwd2 string
	var errMsg string
	var formData url.Values
	currentUser, err := getCurrentUser(c)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		goto end
	}
	if currentUser == nil {
		// should not happen
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpProfile))
	}

	// FIXME this is for demo purpose only
	if demoMode && currentUser.Username == systemUserUsername {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_change_password_system_user_demo")
		goto end
	}

	formData, err = c.FormParams()
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_form_400", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": err.Error()},
		})
		goto end
	}
	currentPwd = strings.TrimSpace(formData.Get("currentPassword"))
	encPwd = encryptPassword(currentUser.Username, currentPwd)
	if encPwd != currentUser.Password {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_password_not_matched")
		goto end
	}
	pwd = strings.TrimSpace(formData.Get("password"))
	pwd2 = strings.TrimSpace(formData.Get("password2"))
	if pwd == "" {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_empty_user_password")
		goto end
	}
	if pwd != pwd2 {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_mismatched_passwords")
		goto end
	}
	currentUser.Password = encryptPassword(currentUser.Username, pwd)
	_, err = userDao.Update(currentUser)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_111", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "change_password_successful"))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_profile", map[string]interface{}{
		"active": "profile",
		"error":  errMsg,
	})
}

/*----------------------------------------------------------------------*/

func actionCpGroupList(c echo.Context) error {
	u := &MyAppUtils{c: c}
	return c.Render(http.StatusOK, namespace+":layout:cp_groups", map[string]interface{}{
		"active":     "groups",
		"userGroups": u.AllUserGroups(),
	})
}

func checkCpCreateGroup(c echo.Context) error {
	if currentUser, err := getCurrentUser(c); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		return errors.New(errMsg)
	} else if currentUser == nil || currentUser.GroupId != systemGroupId {
		// only admin can create groups
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_no_permission")
		return errors.New(errMsg)
	}
	return nil
}

func actionCpCreateGroup(c echo.Context) error {
	if err := checkCpCreateGroup(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}
	formData, _ := c.FormParams()
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active": "groups",
		"form":   formData,
	})
}

func actionCpCreateGroupSubmit(c echo.Context) error {
	if err := checkCpCreateGroup(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}

	var errMsg string
	var err error
	var formData url.Values
	var existingGroup, group *Group

	formData, err = c.FormParams()
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_form_400", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": err.Error()},
		})
		goto end
	}

	group = &Group{
		Id:   strings.ToLower(strings.TrimSpace(formData.Get("id"))),
		Name: strings.TrimSpace(formData.Get("name")),
	}
	if group.Id == "" {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_empty_group_id")
		goto end
	}
	existingGroup, err = groupDao.Get(group.Id)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_301", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": group.Id + "/" + err.Error()},
		})
		goto end
	}
	if existingGroup != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_group_existed", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"group": group.Id},
		})
		goto end
	}
	_, err = groupDao.Create(group.Id, group.Name)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_321", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": group.Id + "/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "create_group_successful", &goyai.LocalizeConfig{
		TemplateData: map[string]interface{}{"group": group.Id},
	}))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active": "groups",
		"form":   formData,
		"error":  errMsg,
	})
}

func checkCpEditGroup(c echo.Context) (*Group, error) {
	gid := c.QueryParam("id")
	if group, err := groupDao.Get(gid); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_301", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": gid + "/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if group == nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_group_not_found", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"group": gid},
		})
		return nil, errors.New(errMsg)
	} else {
		return group, nil
	}
}

func actionCpEditGroup(c echo.Context) error {
	group, err := checkCpEditGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}

	formData := url.Values{}
	formData.Set("id", group.Id)
	formData.Set("name", group.Name)
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active":   "groups",
		"editMode": true,
		"form":     formData,
	})
}

func actionCpEditGroupSubmit(c echo.Context) error {
	group, err := checkCpEditGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}

	var errMsg string
	formData, err := c.FormParams()
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_form_400", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": err.Error()},
		})
		goto end
	}
	group.Name = strings.TrimSpace(formData.Get("name"))
	_, err = groupDao.Update(group)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_311", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": group.Id + "/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "update_group_successful", &goyai.LocalizeConfig{
		TemplateData: map[string]interface{}{"group": group.Id},
	}))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active":   "groups",
		"editMode": true,
		"form":     formData,
		"error":    errMsg,
	})
}

func checkCpDeleteGroup(c echo.Context) (*Group, error) {
	if currentUser, err := getCurrentUser(c); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if currentUser == nil || currentUser.GroupId != systemGroupId {
		// only admin can delete groups
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_no_permission")
		return nil, errors.New(errMsg)
	}
	gid := c.QueryParam("id")
	if group, err := groupDao.Get(gid); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_301", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": gid + "/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if group == nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_group_not_found", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"group": gid},
		})
		return nil, errors.New(errMsg)
	} else if group.Id == systemGroupId {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_delete_system_group", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"group": gid},
		})
		return nil, errors.New(errMsg)
	} else {
		return group, nil
	}
}

func actionCpDeleteGroup(c echo.Context) error {
	group, err := checkCpDeleteGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}

	return c.Render(http.StatusOK, namespace+":layout:cp_delete_group", map[string]interface{}{
		"active":    "groups",
		"userGroup": toGroupModel(c, group),
	})
}

func actionCpDeleteGroupSubmit(c echo.Context) error {
	group, err := checkCpDeleteGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}

	var errMsg string
	_, err = groupDao.Delete(group)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_331", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": group.Id + "/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "delete_group_successful", &goyai.LocalizeConfig{
		TemplateData: map[string]interface{}{"group": group.Id},
	}))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_delete_group", map[string]interface{}{
		"active":    "groups",
		"userGroup": toGroupModel(c, group),
		"error":     errMsg,
	})
}

/*----------------------------------------------------------------------*/

func actionCpUserList(c echo.Context) error {
	u := &MyAppUtils{c: c}
	return c.Render(http.StatusOK, namespace+":layout:cp_users", map[string]interface{}{
		"active": "users",
		"users":  u.AllUsers(),
	})
}

func checkCpCreateUser(c echo.Context) error {
	if currentUser, err := getCurrentUser(c); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		return errors.New(errMsg)
	} else if currentUser == nil || currentUser.GroupId != systemGroupId {
		// only admin can create users
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_no_permission")
		return errors.New(errMsg)
	}
	return nil
}

func actionCpCreateUser(c echo.Context) error {
	if err := checkCpCreateUser(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}
	formData, _ := c.FormParams()
	u := &MyAppUtils{c: c}
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_user", map[string]interface{}{
		"active":     "users",
		"form":       formData,
		"userGroups": u.AllUserGroups(),
	})
}

func actionCpCreateUserSubmit(c echo.Context) error {
	if err := checkCpCreateUser(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+utils.RandomString(4))
	}

	var errMsg string
	var err error
	var formData url.Values
	var existingUser, user *User
	var u = &MyAppUtils{c: c}
	var pwd, pwd2 string

	formData, err = c.FormParams()
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_form_400", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": err.Error()},
		})
		goto end
	}

	user = &User{
		Username: strings.ToLower(strings.TrimSpace(formData.Get("username"))),
		Name:     strings.TrimSpace(formData.Get("name")),
		GroupId:  strings.ToLower(strings.TrimSpace(formData.Get("group"))),
	}
	pwd = strings.TrimSpace(formData.Get("password"))
	pwd2 = strings.TrimSpace(formData.Get("password2"))
	if user.Username == "" {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_empty_user_username")
		goto end
	}
	existingUser, err = userDao.Get(user.Username)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": user.Username + "/" + err.Error()},
		})
		goto end
	}
	if existingUser != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_user_existed", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"user": user.Username},
		})
		goto end
	}
	if pwd == "" {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_empty_user_password")
		goto end
	}
	if pwd != pwd2 {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_mismatched_passwords")
		goto end
	}
	user.Password = encryptPassword(user.Username, pwd)
	_, err = userDao.Create(user.Username, user.Password, user.Name, user.GroupId)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_121", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": user.Username + "/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "create_user_successful", &goyai.LocalizeConfig{
		TemplateData: map[string]interface{}{"user": user.Username},
	}))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_user", map[string]interface{}{
		"active":     "users",
		"form":       formData,
		"userGroups": u.AllUserGroups(),
		"error":      errMsg,
	})
}

func checkCpEditUser(c echo.Context) (*User, error) {
	if currentUser, err := getCurrentUser(c); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if currentUser == nil || currentUser.GroupId != systemGroupId {
		// only admin can edit users
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_no_permission")
		return nil, errors.New(errMsg)
	}
	username := c.QueryParam("u")
	if user, err := userDao.Get(username); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": username + "/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if user == nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_user_not_found", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"user": username},
		})
		return nil, errors.New(errMsg)
	} else if demoMode && username == systemUserUsername {
		// FIXME for demo purpose only
		return nil, errors.New(fmt.Sprintf("Cannot edit system account account [%s]", username))
	} else {
		return user, nil
	}
}

func actionCpEditUser(c echo.Context) error {
	user, err := checkCpEditUser(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
	}

	u := &MyAppUtils{c: c}
	formData := url.Values{}
	formData.Set("username", user.Username)
	formData.Set("name", user.Name)
	formData.Set("group", user.GroupId)
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_user", map[string]interface{}{
		"active":       "users",
		"editMode":     true,
		"form":         formData,
		"userGroups":   u.AllUserGroups(),
		"disableGroup": demoMode && user.Username == systemUserUsername,
	})
}

func actionCpEditUserSubmit(c echo.Context) error {
	user, err := checkCpEditUser(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
	}

	var u = &MyAppUtils{c: c}
	var errMsg string
	var pwd, pwd2 string
	formData, err := c.FormParams()
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_form_400", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": err.Error()},
		})
		goto end
	}
	pwd = strings.TrimSpace(formData.Get("password"))
	pwd2 = strings.TrimSpace(formData.Get("password2"))
	if pwd != "" {
		// to change password: enter new one
		if pwd != pwd2 {
			errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_mismatched_passwords")
			goto end
		}
		user.Password = encryptPassword(user.Username, pwd)
	}
	user.Name = strings.TrimSpace(formData.Get("name"))
	if !demoMode || user.Username != systemUserUsername {
		// do not change group of system admin user
		user.GroupId = strings.ToLower(strings.TrimSpace(formData.Get("group")))
	}
	_, err = userDao.Update(user)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_111", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": user.Username + "/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "update_user_successful", &goyai.LocalizeConfig{
		TemplateData: map[string]interface{}{"user": user.Username},
	}))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_user", map[string]interface{}{
		"active":       "users",
		"editMode":     true,
		"form":         formData,
		"userGroups":   u.AllUserGroups(),
		"error":        errMsg,
		"disableGroup": demoMode && user.Username == systemUserUsername,
	})
}

func checkCpDeleteUser(c echo.Context) (*User, error) {
	if currentUser, err := getCurrentUser(c); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": "current_user/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if currentUser == nil || currentUser.GroupId != systemGroupId {
		// only admin can delete users
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_no_permission")
		return nil, errors.New(errMsg)
	}
	username := c.QueryParam("u")
	if user, err := userDao.Get(username); err != nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_db_101", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": username + "/" + err.Error()},
		})
		return nil, errors.New(errMsg)
	} else if user == nil {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_user_not_found", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"user": username},
		})
		return nil, errors.New(errMsg)
	} else if demoMode && username == systemUserUsername {
		errMsg := myI18n.Localize(getContextString(c, ctxLocale), "error_delete_system_user", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"user": username},
		})
		return nil, errors.New(errMsg)
	} else {
		return user, nil
	}
}

func actionCpDeleteUser(c echo.Context) error {
	user, err := checkCpDeleteUser(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
	}

	return c.Render(http.StatusOK, namespace+":layout:cp_delete_user", map[string]interface{}{
		"active": "users",
		"user":   toUserModel(c, user),
	})
}

func actionCpDeleteUserSubmit(c echo.Context) error {
	user, err := checkCpDeleteUser(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
	}

	var errMsg string
	_, err = userDao.Delete(user)
	if err != nil {
		errMsg = myI18n.Localize(getContextString(c, ctxLocale), "error_db_131", &goyai.LocalizeConfig{
			TemplateData: map[string]interface{}{"err": user.Username + "/" + err.Error()},
		})
		goto end
	}
	addFlashMsg(c, myI18n.Localize(getContextString(c, ctxLocale), "delete_user_successful", &goyai.LocalizeConfig{
		TemplateData: map[string]interface{}{"error": user.Username},
	}))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpUsers)+"?r="+utils.RandomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_delete_user", map[string]interface{}{
		"active": "users",
		"user":   toUserModel(c, user),
		"error":  errMsg,
	})
}
