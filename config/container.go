package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/GabsMeloTI/go_adega/internal/auth"
	"github.com/GabsMeloTI/go_adega/internal/email"
	"github.com/GabsMeloTI/go_adega/internal/metrics"
	"github.com/GabsMeloTI/go_adega/internal/orders"
	"github.com/GabsMeloTI/go_adega/internal/paymentclient"
	"github.com/GabsMeloTI/go_adega/internal/people"
	"github.com/GabsMeloTI/go_adega/internal/products"
	"github.com/GabsMeloTI/go_adega/internal/reports"
	"github.com/GabsMeloTI/go_adega/internal/settings"
	"github.com/GabsMeloTI/go_adega/internal/tracking"
	"github.com/GabsMeloTI/go_adega/internal/uploads"
)

type Container struct {
	Config          Config
	DB              *sql.DB
	ProductHandler  *products.Handler
	OrderHandler    *orders.Handler
	TrackingHandler *tracking.Handler
	SettingsHandler *settings.Handler
	PeopleHandler   *people.Handler
	UploadHandler   *uploads.Handler
	MetricsHandler  *metrics.Handler
	ReportsHandler  *reports.Handler
	AuthHandler     *auth.Handler
}

func NewContainer(cfg Config) (*Container, error) {
	db, err := sql.Open(cfg.DBDriver, cfg.DatabaseDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	productRepo := products.NewRepository(db)
	orderRepo := orders.NewRepository(db)
	trackingRepo := tracking.NewRepository(db)
	settingsRepo := settings.NewRepository(db)
	peopleRepo := people.NewRepository(db)
	metricsRepo := metrics.NewRepository(db)
	reportsRepo := reports.NewRepository(db)
	authRepo := auth.NewRepository(db)
	paymentClient := paymentclient.New(cfg.PaymentServiceURL, http.DefaultClient)
	uploadService, err := uploads.NewService(cfg.GCSBucket, cfg.GCSPublicBaseURL)
	if err != nil {
		return nil, fmt.Errorf("cannot create upload service: %w", err)
	}

	productService := products.NewService(productRepo)
	orderService := orders.NewService(orderRepo, settingsRepo, paymentClient, cfg.PaymentProvider)
	trackingService := tracking.NewService(trackingRepo, orderRepo, settingsRepo)
	settingsService := settings.NewService(settingsRepo)
	emailClient := email.NewSendGridClient(cfg.SendGridAPIKey, cfg.SendGridFromEmail, cfg.SendGridFromName)
	peopleService := people.NewService(peopleRepo, emailClient, cfg.FrontAppURL)
	authService := auth.NewService(authRepo, cfg.AdminEmail, cfg.AdminPassword)

	return &Container{
		Config:          cfg,
		DB:              db,
		ProductHandler:  products.NewHandler(productService),
		OrderHandler:    orders.NewHandler(orderService),
		TrackingHandler: tracking.NewHandler(trackingService),
		SettingsHandler: settings.NewHandler(settingsService),
		PeopleHandler:   people.NewHandler(peopleService),
		UploadHandler:   uploads.NewHandler(uploadService),
		MetricsHandler:  metrics.NewHandler(metricsRepo),
		ReportsHandler:  reports.NewHandler(reportsRepo),
		AuthHandler:     auth.NewHandler(authService),
	}, nil
}

func (c *Container) Close() error {
	if c.DB == nil {
		return nil
	}
	return c.DB.Close()
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("cannot create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migration", "postgres", driver)
	if err != nil {
		return fmt.Errorf("cannot load migrations: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("database migrations: no change")
			return nil
		}
		return fmt.Errorf("cannot run migrations: %w", err)
	}

	log.Println("database migrations: applied")
	return nil
}
