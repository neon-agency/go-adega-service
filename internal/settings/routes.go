package settings

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	settings := g.Group("/settings")
	settings.GET("/store", h.Get)
	settings.PUT("/store", h.Update)
}
