package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetUser(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var userID int64
	var name string

	err := row.Scan(&userID, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	log.Printf("PostgresUserRepository.GetUser success: id=%d name=%s", userID, name)

	return &User{
		ID:   userID,
		Name: name,
	}, nil
}

func (r *PostgresUserRepository) BatchGetUsers(ctx context.Context, ids []int64) ([]*User, error) {
	if len(ids) == 0 {
		return []*User{}, nil
	}

	query := `
		SELECT id, name
		FROM users
		WHERE id = Any($1)
	`
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userMap := make(map[int64]*User, len(ids))

	for rows.Next() {
		var id int64
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		userMap[id] = &User{
			ID:   id,
			Name: name,
		}
	}

	// Checks errors during scanning returning rows, not while query execution
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*User, 0, len(ids))
	for _, id := range ids {
		user, ok := userMap[id]
		if !ok {
			continue
		}

		result = append(result, user)
	}

	return result, nil
}
