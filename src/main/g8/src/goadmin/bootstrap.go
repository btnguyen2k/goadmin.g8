package goadmin

import (
	"github.com/go-akka/configuration"
	"github.com/labstack/echo/v4"
)

/*
IBootstrapper defines an interface for application to hook bootstrapping routines.

Bootstrapper has access to global variables:
- Application configurations via goadmin.AppConfig
- Echo server via goadmin.EchoServer
*/
type IBootstrapper interface {
	Bootstrap(appConfig *configuration.Config, echoServer *echo.Echo) error
}
