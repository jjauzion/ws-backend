package db

import (
	"context"
)

//go:generate mockgen -source=db.go -destination=mock_db/mock_db.go -package=mock_db

// Dbal for DataBase Abstracted Layer
type Dbal interface {
	// newConnection create a connection and store it
	NewConnection(address string) error
	// Ping try to get info from the nodes and return an error if it failed
	Ping() error
	// CreateIndexes initialize needed indexes
	CreateIndexes(ctx context.Context) error

	// User index methods:
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id string) error

	// Task index methods:
	GetNextTask(ctx context.Context) (*Task, error)
	GetTasksByUserID(ctx context.Context, id string) ([]Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status Status) error
	DeleteTaskByID(ctx context.Context, id string) error
	DeleteTasksBysUserID(ctx context.Context, userId string) error
	CreateTask(ctx context.Context, task Task) error
}
