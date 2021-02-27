package db

import (
	"context"
	"github.com/jjauzion/ws-backend/graph/model"
)

type DatabaseHandler interface {
	new() error
	Info() string
	Bootstrap() error

	GetUserByEmail(email string) (model.User, error)
	GetUserByID(id string) (model.User, error)
	CreateUser(user model.User) error

	GetTasksByUserID(ctx context.Context, id string) ([]model.Task, error)
	DeleteTask(ctx context.Context, id string) error
	DeleteUserTasks(ctx context.Context, userId string) error
	CreateTask(task model.Task) error
}
