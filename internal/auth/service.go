package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo          Repository
	adminEmail    string
	adminPassword string
}

func NewService(repo Repository, adminEmail, adminPassword string) *Service {
	return &Service{repo: repo, adminEmail: strings.ToLower(strings.TrimSpace(adminEmail)), adminPassword: adminPassword}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (Session, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)
	req.StoreName = strings.TrimSpace(req.StoreName)
	if req.Name == "" || req.Email == "" || req.Phone == "" || req.StoreName == "" || len(req.Password) < 6 {
		return Session{}, errors.New("preencha nome, e-mail, telefone, nome da loja e senha com no mínimo 6 caracteres")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return Session{}, err
	}
	user, err := s.repo.Create(ctx, req, string(hash))
	if err != nil {
		return Session{}, errors.New("não foi possível cadastrar; verifique se o e-mail já existe")
	}
	return Session{Token: tokenFor(user.ID), User: user}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (Session, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	found, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		if bcrypt.CompareHashAndPassword([]byte(found.PasswordHash), []byte(req.Password)) != nil {
			return Session{}, errors.New("credenciais inválidas")
		}
		return Session{Token: tokenFor(found.ID), User: found.User}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Session{}, err
	}
	if email == s.adminEmail && req.Password == s.adminPassword {
		user := User{
			ID:        "admin",
			Name:      "Administrador",
			Email:     s.adminEmail,
			Role:      "owner",
			StoreName: "Adega",
			Initials:  "AD",
			Color:     "#6E1F2C",
		}
		return Session{Token: tokenFor(user.ID), User: user}, nil
	}
	return Session{}, errors.New("credenciais inválidas")
}

func (s *Service) Me(ctx context.Context, token string) (User, error) {
	id := strings.TrimPrefix(strings.TrimSpace(token), "Bearer ")
	id = strings.TrimPrefix(id, "admin:")
	if id == "" {
		return User{}, errors.New("sessão inválida")
	}
	if id == "admin" {
		return User{
			ID:        "admin",
			Name:      "Administrador",
			Email:     s.adminEmail,
			Role:      "owner",
			StoreName: "Adega",
			Initials:  "AD",
			Color:     "#6E1F2C",
		}, nil
	}
	return s.repo.FindByID(ctx, id)
}

func tokenFor(id string) string {
	return "admin:" + id
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "AD"
	}
	first := []rune(parts[0])
	out := strings.ToUpper(string(first[:1]))
	if len(parts) > 1 {
		last := []rune(parts[len(parts)-1])
		out += strings.ToUpper(string(last[:1]))
	}
	return out
}
