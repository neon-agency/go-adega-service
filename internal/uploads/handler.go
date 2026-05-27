package uploads

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UploadImage(c echo.Context) error {
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "arquivo obrigatório"})
	}
	folder := c.FormValue("folder")
	if folder == "" {
		folder = "products"
	}
	url, err := h.service.Upload(c.Request().Context(), file, header, folder)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"url": url})
}
