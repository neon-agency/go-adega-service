package tracking

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("/orders/:orderId/delivery", h.CreateDelivery)
	g.GET("/deliveries/available", h.ListAvailable)
	g.GET("/drivers/:driverId/deliveries", h.ListByDriver)
	g.POST("/drivers/:driverId/deliveries/:code/claim", h.Claim)
	g.PATCH("/drivers/:driverId/deliveries/:code/location", h.UpdateDriverLocation)
	tracking := g.Group("/tracking")
	tracking.GET("/:code", h.GetByCode)
	tracking.PATCH("/:code/location", h.UpdateLocation)
}
