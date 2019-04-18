package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Validator = newValidator()
	return e
}
