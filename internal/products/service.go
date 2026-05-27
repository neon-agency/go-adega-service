package products

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

func (s *Service) List(ctx context.Context, category string, onlyActive bool) ([]Product, error) {
	return s.repo.List(ctx, category, onlyActive)
}

func (s *Service) Get(ctx context.Context, id string) (Product, error) {
	if strings.TrimSpace(id) == "" {
		return Product{}, errors.New("produto obrigatório")
	}
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Category) == "" || strings.TrimSpace(req.Description) == "" || req.PriceCents < 0 || req.CostCents < 0 {
		return Product{}, errors.New("produto inválido")
	}
	return s.repo.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateProductRequest) (Product, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Category) == "" || strings.TrimSpace(req.Description) == "" || req.PriceCents < 0 || req.CostCents < 0 {
		return Product{}, errors.New("produto inválido")
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("produto obrigatório")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) SetAvailability(ctx context.Context, id string, isActive bool) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("produto obrigatório")
	}
	return s.repo.SetAvailability(ctx, id, isActive)
}

func (s *Service) AddStockMovement(ctx context.Context, productID string, req StockMovementRequest) error {
	if req.Quantity <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}
	switch req.Type {
	case "entry", "sale", "adjustment", "loss":
		return s.repo.AddStockMovement(ctx, productID, req)
	default:
		return errors.New("tipo de movimentação inválido")
	}
}
