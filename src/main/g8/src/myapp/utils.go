package myapp

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get(namespace, c)
	return sess
}
