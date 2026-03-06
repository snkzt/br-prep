package main

import (
	"context"
	"errors"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID   int64
	Name string
}

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	BatchGetUsers(ctx context.Context, ids []int64) ([]*User, error)
}

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[int64]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: map[int64]*User{
			1:  {Name: "Natalie"},
			2:  {Name: "Domingo"},
			11: {Name: "Lucy"},
			99: {Name: "Rasicov"},
		},
	}
}

func (i *InMemoryUserRepository) GetUser(ctx context.Context, id int64) (*User, error) {
	i.mu.RLock()
	user, ok := i.users[id]
	i.mu.RUnlock()

	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Behaviour: preserve input order, ignore unknown IDs, do not deduplicate.
func (i *InMemoryUserRepository) BatchGetUsers(ctx context.Context, ids []int64) ([]*User, error) {
	names := make([]*User, 0, len(ids))

	for _, id := range ids {
		i.mu.RLock()
		user, ok := i.users[id]
		i.mu.RUnlock()
		if ok {
			names = append(names, user)
		}
	}

	return names, nil
}
