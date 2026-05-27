package tracking

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

func (h *Handler) CreateDelivery(c echo.Context) error {
	var req CreateDeliveryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.CreateDelivery(c.Request().Context(), c.Param("orderId"), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) UpdateLocation(c echo.Context) error {
	var req UpdateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	if err := h.service.UpdateLocation(c.Request().Context(), c.Param("code"), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateDriverLocation(c echo.Context) error {
	var req UpdateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	if err := h.service.UpdateDriverLocation(c.Request().Context(), c.Param("driverId"), c.Param("code"), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListByDriver(c echo.Context) error {
	res, err := h.service.ListByDriver(c.Request().Context(), c.Param("driverId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ListAvailable(c echo.Context) error {
	res, err := h.service.ListAvailable(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Claim(c echo.Context) error {
	var req ClaimDeliveryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	if err := h.service.Claim(c.Request().Context(), c.Param("driverId"), c.Param("code"), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetByCode(c echo.Context) error {
	res, err := h.service.GetByCode(c.Request().Context(), c.Param("code"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
