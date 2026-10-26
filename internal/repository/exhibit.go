package repository

import (
	"Museum/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExhibitRepository struct {
	pool *pgxpool.Pool
}

func NewExhibitRepository(pool *pgxpool.Pool) *ExhibitRepository {
	return &ExhibitRepository{pool: pool}
}

func (r *ExhibitRepository) Create(ctx context.Context, e *models.Exhibit) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO exhibits (category, save_by, name, description, 
                      storage_location, status, length, width, height, weight, creation_date)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at, updated_at`,
		e.Category, e.SaveBy, e.Name, e.Description, e.StorageLocation,
		e.Status, e.Length, e.Width, e.Height, e.Weight, e.CreationDate,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *ExhibitRepository) GetById(ctx context.Context, id int64) (*models.Exhibit, error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM exhibits WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Exhibit])
	//var e models.Exhibit
	//err := r.pool.QueryRow(ctx,
	//	`SELECT * FROM exhibits WHERE id = $1`, id).
	//	Scan(&e.ID, &e.Category, &e.SaveBy, &e.Name,
	//		&e.Description, &e.StorageLocation, &e.Status, &e.Length, &e.Width, &e.Height,
	//		&e.CreationDate, &e.CreatedAt, &e.UpdatedAt, &e.Weight)
	//if err != nil {
	//	return nil, err
	//}
	//return &e, nil
}

func (r *ExhibitRepository) List(ctx context.Context, offset, limit int64) ([]models.Exhibit, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM exhibits`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		"SELECT * FROM exhibits ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	//defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Exhibit])

	//items := make([]models.Exhibit, 0, limit)
	//for rows.Next() {
	//	var e models.Exhibit
	//	if err := rows.Scan(&e.ID, &e.Category, &e.SaveBy, &e.Name,
	//		&e.Description, &e.StorageLocation, &e.Status, &e.Length, &e.Width, &e.Height,
	//		&e.CreationDate, &e.CreatedAt, &e.UpdatedAt, &e.Weight); err != nil {
	//		return nil, 0, err
	//	}
	//	items = append(items, e)
	//}
	return items, total, err
}

func (r *ExhibitRepository) Update(ctx context.Context, e *models.Exhibit) error {
	query := `
        UPDATE exhibits SET
            category = $2,
            name = $3,
            description = $4,
            storage_location = $5,
            length = $6,
            width = $7,
            height = $8,
            weight = $9,
            creation_date = $10,
            save_by = $11,
            status = $12
        WHERE id = $1
        RETURNING created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		e.ID,
		e.Category,
		e.Name,
		e.Description,
		e.StorageLocation,
		e.Length,
		e.Width,
		e.Height,
		e.Weight,
		e.CreationDate,
		e.SaveBy,
		e.Status,
	).Scan(&e.CreatedAt, &e.UpdatedAt)
	return err
}

func (r *ExhibitRepository) Exists(ctx context.Context, id int64) bool {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM exhibits
			WHERE id = $1
		);`, id).Scan(&exists)
	return err == nil && exists
}

func (r *ExhibitRepository) MostSuitable(ctx context.Context, vector []float32) ([]models.Exhibit, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT *
		FROM exhibits
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> $1
		LIMIT 5;`, vector)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Exhibit, 0, 5)
	for rows.Next() {
		var e models.Exhibit
		if err := rows.Scan(&e.ID, &e.Category, &e.SaveBy, &e.Name,
			&e.Description, &e.StorageLocation, &e.Status, &e.Length, &e.Width, &e.Height,
			&e.CreationDate, &e.CreatedAt, &e.UpdatedAt, &e.Weight); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}
