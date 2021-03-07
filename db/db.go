package db

import (
	"context"
)

type DatabaseHandler interface {
	new() error
	Info() string
	Bootstrap(ctx context.Context) error

	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id string) error

	GetTasksByUserID(ctx context.Context, id string) ([]Task, error)
	DeleteTask(ctx context.Context, id string) error
	DeleteUserTasks(ctx context.Context, userId string) error
	CreateTask(ctx context.Context, task Task) error
}
