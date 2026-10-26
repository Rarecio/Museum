package repository

import (
	"Museum/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, full_name, role)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		u.Email, u.PasswordHash, u.FullName, u.Role,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, full_name, role, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, full_name, role, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	return r.pool.QueryRow(ctx,
		`UPDATE users SET full_name = $1 WHERE id = $2 RETURNING updated_at`,
		u.FullName, u.ID,
	).Scan(&u.UpdatedAt)
}
