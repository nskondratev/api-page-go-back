package logger

import "github.com/labstack/echo"

type Logger echo.Logger

var l Logger

func New(lg echo.Logger) Logger {
	l = lg
	return l
}
