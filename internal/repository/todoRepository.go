package repository

import (
	"Todo_App/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("todo not found")

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO todos (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, user_id, created_at, updated_at`

	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetTodos(pool *pgxpool.Pool, userID int) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, completed, user_id, created_at, updated_at FROM todos WHERE user_id = $1 ORDER BY id ASC`
	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.UserID,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoByID(pool *pgxpool.Pool, id int, userID int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, completed, user_id, created_at, updated_at FROM todos WHERE id = $1 AND user_id = $2`
	var todo models.Todo
	err := pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, userID int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE todos 
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, completed, user_id, created_at, updated_at`

	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2`
	commandTag, err := pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}