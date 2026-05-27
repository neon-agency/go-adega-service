package people

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/GabsMeloTI/go_adega/internal/email"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo        *Repository
	emailClient email.Client
	frontAppURL string
}

func NewService(repo *Repository, emailClient email.Client, frontAppURL string) *Service {
	return &Service{repo: repo, emailClient: emailClient, frontAppURL: strings.TrimRight(frontAppURL, "/")}
}

func (s *Service) ListDrivers(ctx context.Context) ([]Person, error) {
	return s.repo.ListDrivers(ctx)
}

func (s *Service) LoginDriver(ctx context.Context, req DriverLoginRequest) (Person, error) {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return Person{}, errors.New("e-mail e senha são obrigatórios")
	}
	driver, hash, err := s.repo.LoginDriver(ctx, strings.TrimSpace(req.Email))
	if err != nil {
		return Person{}, err
	}
	if hash == "" || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		return Person{}, errors.New("credenciais inválidas")
	}
	return driver, nil
}

func (s *Service) CreateDriver(ctx context.Context, req UpsertPersonRequest) (Person, error) {
	if req.MaxActiveDeliveries <= 0 {
		req.MaxActiveDeliveries = 1
	}
	if err := validateDriver(req); err != nil {
		return Person{}, err
	}
	tempPassword, err := generateTemporaryPassword()
	if err != nil {
		return Person{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return Person{}, err
	}
	driver, err := s.repo.CreateDriver(ctx, req, string(hash))
	if err != nil {
		return Person{}, err
	}
	if s.emailClient != nil {
		loginURL := fmt.Sprintf("%s/entregador?email=%s", s.frontAppURL, driver.Email)
		if err := s.emailClient.SendDriverWelcome(driver.Email, driver.Name, tempPassword, loginURL); err != nil {
			log.Printf("failed to send driver welcome email: %v", err)
		}
	}
	return driver, nil
}

func (s *Service) UpdateDriver(ctx context.Context, id string, req UpsertPersonRequest) (Person, error) {
	if req.MaxActiveDeliveries <= 0 {
		req.MaxActiveDeliveries = 1
	}
	if err := validateDriver(req); err != nil {
		return Person{}, err
	}
	return s.repo.UpdateDriver(ctx, id, req)
}

func (s *Service) ChangeDriverPassword(ctx context.Context, id string, req DriverChangePasswordRequest) (Person, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(req.CurrentPassword) == "" || len(req.NewPassword) < 6 {
		return Person{}, errors.New("senha inválida")
	}
	_, hash, err := s.repo.GetDriverAuthByID(ctx, id)
	if err != nil {
		return Person{}, err
	}
	if hash == "" || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.CurrentPassword)) != nil {
		return Person{}, errors.New("senha atual inválida")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return Person{}, err
	}
	return s.repo.ChangeDriverPassword(ctx, id, string(newHash))
}

func (s *Service) ListEmployees(ctx context.Context) ([]Person, error) {
	return s.repo.ListEmployees(ctx)
}

func (s *Service) CreateEmployee(ctx context.Context, req UpsertPersonRequest) (Person, error) {
	if req.Role == "" {
		req.Role = "attendant"
	}
	if err := validate(req); err != nil {
		return Person{}, err
	}
	return s.repo.CreateEmployee(ctx, req)
}

func (s *Service) UpdateEmployee(ctx context.Context, id string, req UpsertPersonRequest) (Person, error) {
	if req.Role == "" {
		req.Role = "attendant"
	}
	if err := validate(req); err != nil {
		return Person{}, err
	}
	return s.repo.UpdateEmployee(ctx, id, req)
}

func validate(req UpsertPersonRequest) error {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Phone) == "" {
		return errors.New("nome e telefone são obrigatórios")
	}
	return nil
}

func validateDriver(req UpsertPersonRequest) error {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Email) == "" {
		return errors.New("nome, telefone e e-mail são obrigatórios")
	}
	return nil
}

func generateTemporaryPassword() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")[:10], nil
}
