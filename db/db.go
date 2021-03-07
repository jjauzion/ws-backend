package db

import (
	"context"
)

type DatabaseHandler interface {
	new() error
	Info() string
	Bootstrap(ctx context.Context) error

	GetUserByEmail(email string) (User, error)
	GetUserByID(id string) (User, error)
	CreateUser(user User) error

	GetTasksByUserID(ctx context.Context, id string) ([]Task, error)
	DeleteTask(ctx context.Context, id string) error
	DeleteUserTasks(ctx context.Context, userId string) error
	CreateTask(task Task) error
}
