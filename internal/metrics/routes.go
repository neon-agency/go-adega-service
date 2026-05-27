package metrics

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	metrics := g.Group("/metrics")
	metrics.GET("/overview", h.Overview)
}
