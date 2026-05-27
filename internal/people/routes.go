package people

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	drivers := g.Group("/drivers")
	drivers.GET("", h.ListDrivers)
	drivers.POST("/login", h.LoginDriver)
	drivers.POST("", h.CreateDriver)
	drivers.PATCH("/:id/password", h.ChangeDriverPassword)
	drivers.PUT("/:id", h.UpdateDriver)

	employees := g.Group("/employees")
	employees.GET("", h.ListEmployees)
	employees.POST("", h.CreateEmployee)
	employees.PUT("/:id", h.UpdateEmployee)
}
