package uploads

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("/uploads/images", h.UploadImage)
}
