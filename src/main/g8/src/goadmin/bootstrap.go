package goadmin

import (
	"github.com/go-akka/configuration"
	"github.com/labstack/echo/v4"
)

// IBootstrapper defines an interface for application to hook bootstrapping routines.
type IBootstrapper interface {
	Bootstrap(appConfig *configuration.Config, echoServer *echo.Echo) error
}
