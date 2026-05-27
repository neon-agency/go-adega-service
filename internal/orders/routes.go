package orders

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	orders := g.Group("/orders")
	orders.GET("", h.List)
	orders.GET("/:id", h.Get)
	orders.POST("", h.Create)
	orders.PATCH("/:id/status", h.UpdateStatus)
	orders.DELETE("/:id", h.Cancel)
}
