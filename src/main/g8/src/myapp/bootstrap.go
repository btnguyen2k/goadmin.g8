/*
Package myapp contains application's source code.
*/
package myapp

import (
	"github.com/go-akka/configuration"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"main/src/goadmin"
	"main/src/i18n"
	"net/http"
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
)

const (
	namespace      = "myapp"
	actionNameHome = "home"
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

	// register a custom namespace-scope template renderer
	goadmin.EchoRegisterRenderer(namespace, newTemplateRenderer("./views/myapp", ".html"))

	e.GET("/", actionHome).Name = actionNameHome

	return nil
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

	// sess := getSession(c)
	// flash := sess.Flashes()
	// sess.Save(c.Request(), c.Response())

	// add global data/methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["cdn_mode"] = cdnMode
		viewContext["static"] = myStaticPath
		viewContext["i18n"] = myI18n
		viewContext["reverse"] = c.Echo().Reverse
		viewContext["appInfo"] = goadmin.AppConfig.GetConfig("app")
		// if len(flash) > 0 {
		// 	viewContext["flash"] = flash[0].(string)
		// }
		// uid := c.Get(sessionMyUid)
		// if uid != nil {
		// 	viewContext["uid"] = uid
		// }
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
	return tpl.ExecuteTemplate(w, tokens[0]+".html", data)
}

/*----------------------------------------------------------------------*/

func actionHome(c echo.Context) error {
	err := c.Render(http.StatusOK, namespace+":landing", nil)
	if err != nil {
		log.Error(err)
	}
	return err
}
