/*
Package myapp contains application's source code.
*/
package myapp

import (
	"errors"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/prom"
	"github.com/go-akka/configuration"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"main/src/goadmin"
	"main/src/i18n"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type MyBootstrapper struct {
	name string
}

var (
	Bootstrapper = &MyBootstrapper{name: "myapp"}
	cdnMode      = false
	myStaticPath = "/static"
	myI18n       *i18n.I18n
	sqlc         *prom.SqlConnect
	groupDao     GroupDao
	userDao      UserDao
)

const (
	namespace = "myapp"

	ctxCurrentUser = "usr"
	sessionMyUid   = "uid"

	actionNameHome          = "home"
	actionNameCpLogin       = "cp_login"
	actionNameCpLoginSubmit = "cp_login_submit"
	actionNameCpLogout      = "cp_logout"
	actionNameCpDashboard   = "cp_dashboard"

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

/*
Bootstrap implements goadmin.IBootstrapper.Bootstrap

Bootstrapper usually does:
- register URI routing
- other initializing work (e.g. creating DAO, initializing database, etc)
*/
func (b *MyBootstrapper) Bootstrap(conf *configuration.Config, e *echo.Echo) error {
	cdnMode = conf.GetBoolean(goadmin.ConfKeyCdnMode, false)

	myStaticPath = "/static_v" + conf.GetString("app.version", "")
	e.Static(myStaticPath, "public")

	myI18n = i18n.NewI18n("./config/i18n_" + namespace)

	initDaos()

	// register a custom namespace-scope template renderer
	goadmin.EchoRegisterRenderer(namespace, newTemplateRenderer("./views/myapp", ".html"))

	e.GET("/", actionHome).Name = actionNameHome

	e.GET("/cp/login", actionCpLogin).Name = actionNameCpLogin
	e.POST("/cp/login", actionCpLoginSubmit).Name = actionNameCpLoginSubmit
	// e.GET("/cp/logout", actionCpLogout).Name = actionNameCpLogout
	e.GET("/cp", actionCpDashboard, middlewareRequiredAuth).Name = actionNameCpDashboard

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
	sqlc = newSqliteConnection("./data/sqlite", namespace)
	sqliteInitTableGroup(sqlc, tableGroup)
	sqliteInitTableUser(sqlc, tableUser)

	groupDao = newGroupDaoSqlite(sqlc, tableGroup)
	systemGroup, err := groupDao.Get(SystemGroupId)
	if err != nil {
		panic("error while getting group [" + SystemGroupId + "]: " + err.Error())
	}
	if systemGroup == nil {
		log.Printf("System group [%s] not found, creating one...", SystemGroupId)
		result, err := groupDao.Create(SystemGroupId, "System User Group")
		if err != nil {
			panic("error while creating group [" + SystemGroupId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create group [%s]", SystemGroupId)
		}
	}

	userDao = newUserDaoSqlite(sqlc, tableUser)
	adminUser, err := userDao.Get(AdminUserUsernname)
	if err != nil {
		panic("error while getting user [" + AdminUserUsernname + "]: " + err.Error())
	}
	if adminUser == nil {
		pwd := "s3cr3t"
		log.Printf("Admin user [%s] not found, creating one with password [%s]...", AdminUserUsernname, pwd)
		result, err := userDao.Create(AdminUserUsernname, encryptPassword(AdminUserUsernname, pwd), AdminUserUsernname, SystemGroupId)
		if err != nil {
			panic("error while creating user [" + AdminUserUsernname + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create user [%s]", AdminUserUsernname)
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

/*
Render renders a template document

	- name is list of template names, separated by colon (e.g. <template-name-1>[:<template-name-2>[:<template-name-3>...]])
*/
func (r *myRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
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
				viewContext["user"] = toUserModel(c, &usr)
			case *User:
				viewContext["user"] = toUserModel(c, u.(*User))
			}
		}
	}

	tpl := r.templates[name]
	tokens := strings.Split(name, ":")
	if tpl == nil {
		var files []string
		for _, v := range tokens {
			files = append(files, r.directory+"/"+v+r.templateFileSuffix)
		}
		tpl = template.Must(template.New(name).ParseFiles(files...))
		r.templates[name] = tpl
	}
	// first template-name should be "master" template, and its name is prefixed with ".html"
	return tpl.ExecuteTemplate(w, tokens[0]+".html", data)
}

/*----------------------------------------------------------------------*/
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
	return c.Render(http.StatusOK, namespace+":login", nil)
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
		errMsg = myI18n.Text("error_form_400", err.Error())
		goto end
	}
	username = formData.Get(formFieldUsername)
	user, err = userDao.Get(username)
	if err != nil {
		errMsg = myI18n.Text("error_db_001", err.Error())
		goto end
	}
	if user == nil {
		errMsg = myI18n.Text("error_user_not_found", username)
		goto end
	}
	password = formData.Get(formFieldPassword)
	encPassword = encryptPassword(user.Username, password)
	if encPassword != user.Password {
		errMsg = myI18n.Text("error_login_failed")
		goto end
	}

	// login successful
	setSessionValue(c, sessionMyUid, user.Username)
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpDashboard))
end:
	return c.Render(http.StatusOK, namespace+":login", map[string]interface{}{
		"form":  formData,
		"error": errMsg,
	})
}

func actionCpDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, namespace+":layout:cp_dashboard", map[string]interface{}{
		"active":  "dashboard",
		"osUtils": &OsUtils{},
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
		return errors.New(myI18n.Text("error_db_101", "current_user/"+err.Error()))
	} else if currentUser == nil || currentUser.GroupId != SystemGroupId {
		return errors.New(myI18n.Text("error_no_permission"))
	}
	return nil
}

func actionCpCreateGroup(c echo.Context) error {
	if err := checkCpCreateGroup(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	var errMsg string
	var err error
	var formData url.Values
	var existingGroup, group *Group

	formData, err = c.FormParams()
	if err != nil {
		errMsg = myI18n.Text("error_form_400", err.Error())
		goto end
	}

	group = &Group{
		Id:   strings.ToLower(strings.TrimSpace(formData.Get("id"))),
		Name: strings.TrimSpace(formData.Get("name")),
	}
	if group.Id == "" {
		errMsg = myI18n.Text("error_empty_group_id")
		goto end
	}
	existingGroup, err = groupDao.Get(group.Id)
	if err != nil {
		errMsg = myI18n.Text("error_db_101", group.Id+"/"+err.Error())
		goto end
	}
	if existingGroup != nil {
		errMsg = myI18n.Text("error_group_existed", group.Id)
		goto end
	}
	_, err = groupDao.Create(group.Id, group.Name)
	if err != nil {
		errMsg = myI18n.Text("error_create_group", group.Id, err.Error())
		goto end
	}
	addFlashMsg(c, myI18n.Text("create_group_successful", group.Id))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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
		return nil, errors.New(myI18n.Text("error_db_101", "current/"+err.Error()))
	} else if group == nil {
		return nil, errors.New(myI18n.Text("error_group_not_found", gid))
	} else {
		return group, nil
	}
}

func actionCpEditGroup(c echo.Context) error {
	group, err := checkCpEditGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	var errMsg string
	formData, err := c.FormParams()
	if err != nil {
		errMsg = myI18n.Text("error_form_400", err.Error())
		goto end
	}
	group.Name = strings.TrimSpace(formData.Get("name"))
	_, err = groupDao.Update(group)
	if err != nil {
		errMsg = myI18n.Text("error_update_group", group.Id, err.Error())
		goto end
	}
	addFlashMsg(c, myI18n.Text("update_group_successful", group.Id))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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
		return nil, errors.New(myI18n.Text("error_db_101", "current_user/"+err.Error()))
	} else if currentUser == nil || currentUser.GroupId != SystemGroupId {
		return nil, errors.New(myI18n.Text("error_no_permission"))
	}
	gid := c.QueryParam("id")
	if group, err := groupDao.Get(gid); err != nil {
		return nil, errors.New(myI18n.Text("error_db_101", "current/"+err.Error()))
	} else if group == nil {
		return nil, errors.New(myI18n.Text("error_group_not_found", gid))
	} else if group.Id == SystemGroupId {
		return nil, errors.New(myI18n.Text("error_delete_system_group", gid))
	} else {
		return group, nil
	}
}

func actionCpDeleteGroup(c echo.Context) error {
	group, err := checkCpDeleteGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	var errMsg string
	_, err = groupDao.Delete(group)
	if err != nil {
		errMsg = myI18n.Text("error_delete_group", group.Id, err.Error())
		goto end
	}
	addFlashMsg(c, myI18n.Text("delete_group_successful", group.Id))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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
		return errors.New(myI18n.Text("error_db_101", "current_user/"+err.Error()))
	} else if currentUser == nil || currentUser.GroupId != SystemGroupId {
		return errors.New(myI18n.Text("error_no_permission"))
	}
	return nil
}

func actionCpCreateUser(c echo.Context) error {
	if err := checkCpCreateGroup(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}
	formData, _ := c.FormParams()
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active": "groups",
		"form":   formData,
	})
}

func actionCpCreateUserSubmit(c echo.Context) error {
	if err := checkCpCreateGroup(c); err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	var errMsg string
	var err error
	var formData url.Values
	var existingGroup, group *Group

	formData, err = c.FormParams()
	if err != nil {
		errMsg = myI18n.Text("error_form_400", err.Error())
		goto end
	}

	group = &Group{
		Id:   strings.ToLower(strings.TrimSpace(formData.Get("id"))),
		Name: strings.TrimSpace(formData.Get("name")),
	}
	if group.Id == "" {
		errMsg = myI18n.Text("error_empty_group_id")
		goto end
	}
	existingGroup, err = groupDao.Get(group.Id)
	if err != nil {
		errMsg = myI18n.Text("error_db_101", group.Id+"/"+err.Error())
		goto end
	}
	if existingGroup != nil {
		errMsg = myI18n.Text("error_group_existed", group.Id)
		goto end
	}
	_, err = groupDao.Create(group.Id, group.Name)
	if err != nil {
		errMsg = myI18n.Text("error_create_group", group.Id, err.Error())
		goto end
	}
	addFlashMsg(c, myI18n.Text("create_group_successful", group.Id))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active": "groups",
		"form":   formData,
		"error":  errMsg,
	})
}

func checkCpEditUser(c echo.Context) (*Group, error) {
	gid := c.QueryParam("id")
	if group, err := groupDao.Get(gid); err != nil {
		return nil, errors.New(myI18n.Text("error_db_101", "current/"+err.Error()))
	} else if group == nil {
		return nil, errors.New(myI18n.Text("error_group_not_found", gid))
	} else {
		return group, nil
	}
}

func actionCpEditUser(c echo.Context) error {
	group, err := checkCpEditGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
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

func actionCpEditUserSubmit(c echo.Context) error {
	group, err := checkCpEditGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	var errMsg string
	formData, err := c.FormParams()
	if err != nil {
		errMsg = myI18n.Text("error_form_400", err.Error())
		goto end
	}
	group.Name = strings.TrimSpace(formData.Get("name"))
	_, err = groupDao.Update(group)
	if err != nil {
		errMsg = myI18n.Text("error_update_group", group.Id, err.Error())
		goto end
	}
	addFlashMsg(c, myI18n.Text("update_group_successful", group.Id))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_create_edit_group", map[string]interface{}{
		"active":   "groups",
		"editMode": true,
		"form":     formData,
		"error":    errMsg,
	})
}

func checkCpDeleteUser(c echo.Context) (*Group, error) {
	if currentUser, err := getCurrentUser(c); err != nil {
		return nil, errors.New(myI18n.Text("error_db_101", "current_user/"+err.Error()))
	} else if currentUser == nil || currentUser.GroupId != SystemGroupId {
		return nil, errors.New(myI18n.Text("error_no_permission"))
	}
	gid := c.QueryParam("id")
	if group, err := groupDao.Get(gid); err != nil {
		return nil, errors.New(myI18n.Text("error_db_101", "current/"+err.Error()))
	} else if group == nil {
		return nil, errors.New(myI18n.Text("error_group_not_found", gid))
	} else if group.Id == SystemGroupId {
		return nil, errors.New(myI18n.Text("error_delete_system_group", gid))
	} else {
		return group, nil
	}
}

func actionCpDeleteUser(c echo.Context) error {
	group, err := checkCpDeleteGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	return c.Render(http.StatusOK, namespace+":layout:cp_delete_group", map[string]interface{}{
		"active":    "groups",
		"userGroup": toGroupModel(c, group),
	})
}

func actionCpDeleteUserSubmit(c echo.Context) error {
	group, err := checkCpDeleteGroup(c)
	if err != nil {
		addFlashMsg(c, flashPrefixWarning+err.Error())
		return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
	}

	var errMsg string
	_, err = groupDao.Delete(group)
	if err != nil {
		errMsg = myI18n.Text("error_delete_group", group.Id, err.Error())
		goto end
	}
	addFlashMsg(c, myI18n.Text("delete_group_successful", group.Id))
	return c.Redirect(http.StatusFound, c.Echo().Reverse(actionNameCpGroups)+"?r="+randomString(4))
end:
	return c.Render(http.StatusOK, namespace+":layout:cp_delete_group", map[string]interface{}{
		"active":    "groups",
		"userGroup": toGroupModel(c, group),
		"error":     errMsg,
	})
}
