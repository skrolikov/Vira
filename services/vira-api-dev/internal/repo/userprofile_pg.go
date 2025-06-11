package repo

import (
	"context"
	"database/sql"
	"vira-api-dev/internal/types"
)

type PostgresUserProfileRepo struct {
	db *sql.DB
}

func NewUserProfileRepo(db *sql.DB) *PostgresUserProfileRepo {
	return &PostgresUserProfileRepo{db: db}
}

func (r *PostgresUserProfileRepo) Exists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id=$1)", userID).Scan(&exists)
	return exists, err
}

func (r *PostgresUserProfileRepo) Create(ctx context.Context, profile types.UserProfile) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO user_profiles (user_id, city, joined_at) VALUES ($1,$2,$3)",
		profile.UserID, profile.City, profile.JoinedAt)
	return err
}

func (r *PostgresUserProfileRepo) GetByUserID(ctx context.Context, userID string) (types.UserProfile, error) {
	var p types.UserProfile
	err := r.db.QueryRowContext(ctx,
		"SELECT user_id, city, joined_at FROM user_profiles WHERE user_id=$1", userID,
	).Scan(&p.UserID, &p.City, &p.JoinedAt)
	if err == sql.ErrNoRows {
		return types.UserProfile{}, types.ErrProfileNotFound
	}
	return p, err
}
