package myapp

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"strings"
)

func getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get(namespace, c)
	return sess
}

func setSessionValue(c echo.Context, key string, value interface{}) {
	sess := getSession(c)
	if value == nil {
		delete(sess.Values, key)
	} else {
		sess.Values[key] = value
	}
	sess.Save(c.Request(), c.Response())
}

func encryptPassword(username, rawPassword string) string {
	saltAndPwd := username + "." + rawPassword
	out := sha1.Sum([]byte(saltAndPwd))
	return strings.ToLower(hex.EncodeToString(out[:]))
}
