package goadmin

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	TemplateRenderer *GoadminRenderer
)

// EchoSetDefaultRenderer registers the default template renderer instance for the whole application
func EchoSetDefaultRenderer(renderer echo.Renderer) {
	TemplateRenderer.defaultRenderer = renderer
}

// EchoRegisterRenderer registers a namespace-scope renderers
func EchoRegisterRenderer(namespace string, renderer echo.Renderer) {
	if TemplateRenderer.renderers == nil {
		TemplateRenderer.renderers = make(map[string]echo.Renderer)
	}
	TemplateRenderer.renderers[namespace] = renderer
}

func newGoadminRenderer() *GoadminRenderer {
	return &GoadminRenderer{defaultRenderer: &jsonRenderer{}, renderers: make(map[string]echo.Renderer)}
}

// Ref: https://echo.labstack.com/guide/templates

// GoadminRenderer is a routing template renderer that delegates rendering tasks to underlying renderers
type GoadminRenderer struct {
	defaultRenderer echo.Renderer
	renderers       map[string]echo.Renderer
}

// Render implements Renderer.Render.
//
// 'name' must follow the format <namespace>:<template-name>. If there is a renderer associated with the namespace,
// it will be used for rendering (template-name is passed to the renderer via the 'name' parameter); otherwise, the
// default renderer is used.
func (r *GoadminRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	namespaceAndTplname := strings.SplitN(name, ":", 2)
	if renderer, ok := r.renderers[namespaceAndTplname[0]]; ok {
		return renderer.Render(w, namespaceAndTplname[1], data, c)
	}
	return r.defaultRenderer.Render(w, name, data, c)
}

type jsonRenderer struct {
}

// Render implements Renderer.Render
func (r *jsonRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	v := reflect.ValueOf(data)
	if data == nil || v.IsNil() {
		data = make(map[string]interface{})
	}
	if m, ok := data.(map[string]interface{}); ok {
		m["_name_"] = name
		m["_desc_"] = "generated by default 'jsonRenderer'"
		m["_app_"] = AppConfig.GetString("app.name") + " v" + AppConfig.GetString("app.version")
	}
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(js)
	return err
}
