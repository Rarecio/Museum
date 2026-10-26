package repository

import (
	"Museum/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO category (name) VALUES ($1)`, c.Name)
	return err
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.pool.Query(ctx, "SELECT name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Category, 0, 10)
	for rows.Next() {
		var e models.Category
		if err := rows.Scan(&e.Name); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}
