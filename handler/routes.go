package handler

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"strings"
)

func (h *Handler) Register(rg *echo.Group, bg *echo.Group) {
	// Base group
	bg.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: SPASkipper,
		Root:    "public",
		HTML5:   true,
		Browse:  false,
	}))

	// Events routes
	event := rg.Group("/events")
	event.GET("", h.ListEvents)
	event.POST("", h.CreateEvent)
	event.GET("/:id", h.GetEvent)
	event.POST("/:id", h.UpdateEvent)
	event.DELETE("/:id", h.DeleteEvent)

	// Pages routes
	page := rg.Group("/pages")
	page.GET("", h.ListPages)
	page.POST("", h.CreatePage)
	page.GET("/:id", h.GetPage)
	page.POST("/:id", h.UpdatePage)
	page.DELETE("/:id", h.DeletePage)

	// WebSocket route
	rg.GET("/ws", h.HandleWs)

	// GraphQL route
	rg.POST("/graphql", h.handleGraphQLQuery)
}

func SPASkipper(c echo.Context) bool {
	c.Logger().Debugf("SPASkipper call. Request URI: %s", c.Request().RequestURI)
	return strings.Contains(c.Request().RequestURI, "/api/")
}
