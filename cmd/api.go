package cmd

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GabsMeloTI/go_adega/config"
	"github.com/GabsMeloTI/go_adega/internal/auth"
	"github.com/GabsMeloTI/go_adega/internal/docs"
	"github.com/GabsMeloTI/go_adega/internal/metrics"
	"github.com/GabsMeloTI/go_adega/internal/orders"
	"github.com/GabsMeloTI/go_adega/internal/people"
	"github.com/GabsMeloTI/go_adega/internal/products"
	"github.com/GabsMeloTI/go_adega/internal/reports"
	"github.com/GabsMeloTI/go_adega/internal/settings"
	"github.com/GabsMeloTI/go_adega/internal/tracking"
	"github.com/GabsMeloTI/go_adega/internal/uploads"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func StartAPI(ctx context.Context, container *config.Container) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(echoMiddleware.RequestID())
	docs.RegisterRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "ok",
			"service":   "go_adega",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	v1 := e.Group("/api/v1")
	products.RegisterRoutes(v1, container.ProductHandler)
	orders.RegisterRoutes(v1, container.OrderHandler)
	tracking.RegisterRoutes(v1, container.TrackingHandler)
	settings.RegisterRoutes(v1, container.SettingsHandler)
	people.RegisterRoutes(v1, container.PeopleHandler)
	uploads.RegisterRoutes(v1, container.UploadHandler)
	metrics.RegisterRoutes(v1, container.MetricsHandler)
	reports.RegisterRoutes(v1, container.ReportsHandler)
	auth.RegisterRoutes(v1, container.AuthHandler)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Printf("api shutdown failed: %v", err)
		}
	}()

	log.Printf("go_adega api listening on :%s", container.Config.ServerPort)
	e.Logger.Fatal(e.Start(":" + container.Config.ServerPort))
}
