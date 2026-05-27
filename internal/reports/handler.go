package reports

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Get(c echo.Context) error {
	data, err := h.repo.Get(c.Request().Context(), c.QueryParam("period"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}
