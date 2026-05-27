package auth

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

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	session, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "credenciais inválidas"})
	}
	return c.JSON(http.StatusOK, session)
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}
	session, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, session)
}

func (h *Handler) Me(c echo.Context) error {
	user, err := h.service.Me(c.Request().Context(), c.Request().Header.Get(echo.HeaderAuthorization))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "sessão inválida"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Logout(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
