package products

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	products := g.Group("/products")
	products.GET("", h.List)
	products.GET("/:id", h.Get)
	products.POST("", h.Create)
	products.PUT("/:id", h.Update)
	products.DELETE("/:id", h.Delete)
	products.PATCH("/:id/availability", h.SetAvailability)
	products.POST("/:id/stock-movements", h.AddStockMovement)
}
