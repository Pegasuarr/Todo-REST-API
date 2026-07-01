package repository

import (
	"Todo_App/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, username, passwordHash string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash, created_at, updated_at`

	var user models.User
	err := pool.QueryRow(ctx, query, username, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(pool *pgxpool.Pool, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1`
	var user models.User
	err := pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = $1`
	var user models.User
	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
