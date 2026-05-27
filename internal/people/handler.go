package people

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

func (h *Handler) ListDrivers(c echo.Context) error {
	res, err := h.service.ListDrivers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) LoginDriver(c echo.Context) error {
	var req DriverLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.LoginDriver(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "motoboy não encontrado ou inativo"})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ChangeDriverPassword(c echo.Context) error {
	var req DriverChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.ChangeDriverPassword(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateDriver(c echo.Context) error {
	var req UpsertPersonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.CreateDriver(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) UpdateDriver(c echo.Context) error {
	var req UpsertPersonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.UpdateDriver(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ListEmployees(c echo.Context) error {
	res, err := h.service.ListEmployees(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateEmployee(c echo.Context) error {
	var req UpsertPersonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.CreateEmployee(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) UpdateEmployee(c echo.Context) error {
	var req UpsertPersonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	res, err := h.service.UpdateEmployee(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
