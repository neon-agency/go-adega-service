package tracking

import (
	"context"
	"errors"
	"math"

	"github.com/GabsMeloTI/go_adega/internal/orders"
	"github.com/GabsMeloTI/go_adega/internal/settings"
)

type Service struct {
	repo         Repository
	orderRepo    orders.Repository
	settingsRepo settings.Repository
}

func NewService(repo Repository, orderRepo orders.Repository, settingsRepo settings.Repository) *Service {
	return &Service{repo: repo, orderRepo: orderRepo, settingsRepo: settingsRepo}
}

func (s *Service) CreateDelivery(ctx context.Context, orderID string, req CreateDeliveryRequest) (Tracking, error) {
	if orderID == "" {
		return Tracking{}, errors.New("pedido obrigatório")
	}
	if req.DriverID != "" {
		if err := s.orderRepo.UpdateStatus(ctx, orderID, "out_for_delivery"); err != nil {
			return Tracking{}, err
		}
	} else {
		if err := s.orderRepo.UpdateStatus(ctx, orderID, "separating"); err != nil {
			return Tracking{}, err
		}
	}
	return s.repo.CreateDelivery(ctx, orderID, req)
}

func (s *Service) UpdateLocation(ctx context.Context, trackingCode string, req UpdateLocationRequest) error {
	if req.Status == "" {
		req.Status = "out_for_delivery"
	}
	tracking, err := s.repo.GetByCode(ctx, trackingCode)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateLocation(ctx, trackingCode, req); err != nil {
		return err
	}
	if req.Status == "out_for_delivery" || req.Status == "delivered" {
		return s.orderRepo.UpdateStatus(ctx, tracking.OrderID, req.Status)
	}
	return nil
}

func (s *Service) UpdateDriverLocation(ctx context.Context, driverID string, trackingCode string, req UpdateLocationRequest) error {
	if driverID == "" {
		return errors.New("motoboy obrigatório")
	}
	if req.Status == "" {
		req.Status = "out_for_delivery"
	}
	tracking, err := s.repo.GetByCode(ctx, trackingCode)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateDriverLocation(ctx, driverID, trackingCode, req); err != nil {
		return err
	}
	if req.Status == "out_for_delivery" || req.Status == "delivered" {
		return s.orderRepo.UpdateStatus(ctx, tracking.OrderID, req.Status)
	}
	return nil
}

func (s *Service) GetByCode(ctx context.Context, trackingCode string) (Tracking, error) {
	return s.repo.GetByCode(ctx, trackingCode)
}

func (s *Service) ListByDriver(ctx context.Context, driverID string) ([]Tracking, error) {
	if driverID == "" {
		return nil, errors.New("motoboy obrigatório")
	}
	return s.repo.ListByDriver(ctx, driverID)
}

func (s *Service) ListAvailable(ctx context.Context) ([]Tracking, error) {
	return s.repo.ListAvailable(ctx)
}

func (s *Service) Claim(ctx context.Context, driverID string, trackingCode string, req ClaimDeliveryRequest) error {
	if driverID == "" || trackingCode == "" {
		return errors.New("motoboy e entrega são obrigatórios")
	}
	store, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return err
	}
	if store.Latitude == 0 || store.Longitude == 0 {
		return errors.New("localização da loja não configurada")
	}
	if distanceMeters(store.Latitude, store.Longitude, req.Latitude, req.Longitude) > float64(store.DriverPickupRadiusMeters) {
		return errors.New("você precisa estar no estabelecimento para pegar pedidos")
	}
	active, max, err := s.repo.CountActiveByDriver(ctx, driverID)
	if err != nil {
		return err
	}
	if active >= max {
		return errors.New("limite de pedidos ativos atingido")
	}
	if err := s.repo.Claim(ctx, driverID, trackingCode); err != nil {
		return err
	}
	tracking, err := s.repo.GetByCode(ctx, trackingCode)
	if err == nil {
		return s.orderRepo.UpdateStatus(ctx, tracking.OrderID, "out_for_delivery")
	}
	return nil
}

func distanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earth = 6371000
	toRad := func(v float64) float64 { return v * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	return earth * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
