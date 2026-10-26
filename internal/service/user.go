package service

import (
	"Museum/internal/auth"
	"Museum/internal/models"
	"Museum/internal/repository"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmailExists = errors.New("email already exists")

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName string, role models.UserRole) (*models.User, string, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, "", ErrEmailExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	u := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		Role:         role,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, "", err
	}
	token, err := auth.NewToken(u.ID, u.Email, string(u.Role), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", err
	}
	token, err := auth.NewToken(u.ID, u.Email, string(u.Role), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, uid int64, fullName string) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	u.FullName = fullName
	err = s.repo.Update(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return u, nil
}
