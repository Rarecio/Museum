package service

import (
	"Museum/internal/models"
	"Museum/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"time"
)

type ExhibitService struct {
	repo *repository.ExhibitRepository
}

func NewExhibitService(repo *repository.ExhibitRepository) *ExhibitService {
	return &ExhibitService{repo: repo}
}

func (s *ExhibitService) Create(ctx context.Context, category string, saveBy *int64, name, description,
	storageLocation string, length, width, height, weight int32, creationDate time.Time) (*models.Exhibit, error) {
	if category == "" {
		category = "без категории"
	}
	e := &models.Exhibit{
		Category:        category,
		SaveBy:          saveBy,
		Name:            name,
		Description:     description,
		StorageLocation: storageLocation,
		Length:          length,
		Width:           width,
		Height:          height,
		Weight:          weight,
		CreationDate:    creationDate,
		Status:          models.ApproveStatus,
	}
	slog.Debug(e.Category)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("failed to create exhibit: %w", err)
	}
	return e, nil
}

func (s *ExhibitService) GetById(ctx context.Context, id int64) (*models.Exhibit, error) {
	e, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get exhibit by id: %w", err)
	}
	return e, nil
}

func (s *ExhibitService) List(ctx context.Context, page, size int64) (*models.Page[models.Exhibit], error) {
	list, total, err := s.repo.List(ctx, page*size, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of exhibits: %w", err)
	}
	totalPages := int64(math.Ceil(float64(total / size)))
	return models.NewPage(list, total, totalPages, page), nil
}

func (s *ExhibitService) Update(ctx context.Context, uid int64, e *models.Exhibit) error {
	ex, err := s.repo.GetById(ctx, e.ID)
	if err != nil {
		return ErrNotFound
	}
	e.Status = models.ApproveStatus
	e.SaveBy = &uid
	if e.StorageLocation == "" {
		e.StorageLocation = ex.StorageLocation
	} else {
		err = os.Remove(ex.StorageLocation)
		if err != nil {
			slog.Debug("failed to remove exhibit storage location", "location", ex.StorageLocation, "err", err)
		}
	}
	err = s.repo.Update(ctx, e)
	if err != nil {
		return fmt.Errorf("failed to update exhibit: %w", err)
	}
	return nil
}
