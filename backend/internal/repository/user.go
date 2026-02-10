package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vshakitskiy/cis/internal/model"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (email, password, name, role)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		user.Email, user.Password, user.Name, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users WHERE email = $1
	`

	rows, _ := r.pool.Query(ctx, query, email)
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.User])
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
    FROM users WHERE id = $1
	`

	rows, _ := r.pool.Query(ctx, query, id)
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.User])
	if err != nil {
		return nil, err
	}

	return u, nil
}
