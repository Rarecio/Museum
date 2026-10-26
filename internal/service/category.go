package service

import (
	"Museum/internal/models"
	"Museum/internal/repository"
	"context"
	"log/slog"
	"strings"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, name string) (*models.Category, error) {
	c := &models.Category{
		Name: strings.ToLower(strings.TrimSpace(name)),
	}
	err := s.repo.Create(ctx, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	slog.Info("CategoryService List", "list", list)
	return list, nil
}
