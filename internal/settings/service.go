package settings

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context) (StoreSettings, error) {
	return s.repo.Get(ctx)
}

func (s *Service) Update(ctx context.Context, req UpdateStoreSettingsRequest) (StoreSettings, error) {
	if strings.TrimSpace(req.StoreName) == "" || strings.TrimSpace(req.Phone) == "" {
		return StoreSettings{}, errors.New("nome e telefone da adega são obrigatórios")
	}
	if strings.TrimSpace(req.AddressStreet) == "" || strings.TrimSpace(req.AddressNumber) == "" ||
		strings.TrimSpace(req.AddressNeighborhood) == "" || strings.TrimSpace(req.AddressCity) == "" ||
		strings.TrimSpace(req.AddressState) == "" || strings.TrimSpace(req.AddressZipCode) == "" {
		return StoreSettings{}, errors.New("endereço da adega incompleto")
	}
	if req.DeliveryFeeCents < 0 || req.FreeDeliveryFromCents < 0 || req.MinOrderCents < 0 {
		return StoreSettings{}, errors.New("valores financeiros não podem ser negativos")
	}
	if !req.AllowDelivery && !req.AllowPickup {
		return StoreSettings{}, errors.New("habilite entrega ou retirada")
	}
	if strings.TrimSpace(req.BrandColor) == "" {
		req.BrandColor = "#6E1F2C"
	}
	if req.DeliveryRadiusKM <= 0 {
		req.DeliveryRadiusKM = 5
	}
	if req.AverageDeliveryMin <= 0 {
		req.AverageDeliveryMin = 30
	}
	if req.AverageDeliveryMax <= 0 {
		req.AverageDeliveryMax = 45
	}
	if req.AverageDeliveryMax < req.AverageDeliveryMin {
		return StoreSettings{}, errors.New("tempo máximo de entrega deve ser maior ou igual ao mínimo")
	}
	if req.DriverPickupRadiusMeters <= 0 {
		req.DriverPickupRadiusMeters = 150
	}
	if strings.TrimSpace(req.OpeningTime) == "" || strings.TrimSpace(req.ClosingTime) == "" {
		return StoreSettings{}, errors.New("horário de abertura e fechamento são obrigatórios")
	}
	if !req.AcceptOnlinePix && !req.AcceptOnlineCard && !req.AcceptDeliveryPix && !req.AcceptDeliveryCard && !req.AcceptDeliveryCash {
		return StoreSettings{}, errors.New("habilite pelo menos uma forma de pagamento")
	}
	if len(req.OpeningHours) == 0 {
		req.OpeningHours = []OpeningHour{
			{ID: "seg", Label: "Segunda", Open: true, From: req.OpeningTime, To: req.ClosingTime},
			{ID: "ter", Label: "Terça", Open: true, From: req.OpeningTime, To: req.ClosingTime},
			{ID: "qua", Label: "Quarta", Open: true, From: req.OpeningTime, To: req.ClosingTime},
			{ID: "qui", Label: "Quinta", Open: true, From: req.OpeningTime, To: req.ClosingTime},
			{ID: "sex", Label: "Sexta", Open: true, From: req.OpeningTime, To: req.ClosingTime},
			{ID: "sab", Label: "Sábado", Open: true, From: req.OpeningTime, To: req.ClosingTime},
			{ID: "dom", Label: "Domingo", Open: false, From: req.OpeningTime, To: req.ClosingTime},
		}
	}
	return s.repo.Update(ctx, req)
}
