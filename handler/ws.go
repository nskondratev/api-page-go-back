package handler

import (
	"github.com/labstack/echo"
)

func (h *Handler) HandleWs(c echo.Context) error {
	h.wsHub.ServeWs(c.Response(), c.Request())
	return nil
}
