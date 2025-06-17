package repo

import (
	"context"
	"database/sql"
	_ "embed"
	"vira-api-dev/internal/types"
)

//go:embed queries/user_exists.sql
var queryUserExists string

//go:embed queries/user_insert.sql
var queryUserInsert string

//go:embed queries/user_get_by_id.sql
var queryUserGetByID string

type PostgresUserProfileRepo struct {
	db *sql.DB
}

func NewUserProfileRepo(db *sql.DB) *PostgresUserProfileRepo {
	return &PostgresUserProfileRepo{db: db}
}

func (r *PostgresUserProfileRepo) Exists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, queryUserExists, userID).Scan(&exists)
	return exists, err
}

func (r *PostgresUserProfileRepo) Create(ctx context.Context, profile types.UserProfile) error {
	_, err := r.db.ExecContext(ctx, queryUserInsert, profile.UserID, profile.City, profile.JoinedAt)
	return err
}

func (r *PostgresUserProfileRepo) GetByUserID(ctx context.Context, userID string) (types.UserProfile, error) {
	var p types.UserProfile
	err := r.db.QueryRowContext(ctx, queryUserGetByID, userID).
		Scan(&p.UserID, &p.City, &p.JoinedAt)

	if err == sql.ErrNoRows {
		return types.UserProfile{}, types.ErrProfileNotFound
	}
	return p, err
}
