package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/42-AI/ws-backend/internal/auth"
	"time"

	"github.com/google/uuid"
	"github.com/42-AI/ws-backend/db"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*User, error) {
	dbu, err := r.Auth.UserFromContext(ctx, auth.OptOnlyAdmin)
	if err != nil {
		return nil, err
	}

	r.Log.Debug("create user...", zap.String("user_id", dbu.ID))
	newUser := db.User{
		ID:        uuid.New().String(),
		Admin:     false,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}
	err = r.Dbal.CreateUser(ctx, newUser)
	if err != nil {
		if err == db.ErrTooManyHits {
			return nil, fmt.Errorf("user already exist")
		}
		r.Log.Warn("create user: ", zap.Error(err))
		return nil, err
	}

	return UserFromDBModel(newUser).Ptr(), nil
}

func (r *mutationResolver) CreateTask(ctx context.Context, input NewTask) (*Task, error) {
	dbu, err := r.Auth.UserFromContext(ctx, auth.OptAuthenticatedUser)
	if err != nil {
		return nil, err
	}

	r.Log.Debug("create tasks...", zap.String("user_id", dbu.ID))
	user := UserFromDBModel(dbu)
	newJob := db.Job{
		DockerImage: input.DockerImage,
		Dataset:     *input.Dataset,
		Env:         input.Env,
	}
	newTask := db.Task{
		ID:        uuid.New().String(),
		UserId:    user.ID,
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0, 0),
		EndedAt:   time.Unix(0, 0),
		Status:    db.StatusNotStarted,
		Job:       newJob,
	}
	if err = r.Dbal.CreateTask(ctx, newTask); err != nil {
		return nil, err
	}
	r.Log.Info("task created", zap.String("id", newTask.ID))
	return TaskFromDBModel(newTask).Ptr(), nil
}

func (r *queryResolver) ListTasks(ctx context.Context) ([]*Task, error) {
	user, err := r.Auth.UserFromContext(ctx, auth.OptAuthenticatedUser)
	if err != nil {
		return nil, err
	}

	r.Log.Debug("list tasks...", zap.String("user_id", user.ID))
	r.Log.Debug("user authenticated", zap.String("user_email", user.Email))
	res, err := r.Dbal.GetTasksByUserID(ctx, user.ID)
	if err != nil {
		r.Log.Warn("cannot get tasks", zap.String("user_id", user.ID), zap.Error(err))
		return nil, err
	}

	tasks := Tasks{}
	for _, re := range res {
		tasks = append(tasks, TaskFromDBModel(re).Ptr())
	}

	r.Log.Info("list tasks success",
		zap.Int("tasks found", len(tasks)),
		zap.String("user_email", user.Email))
	r.Log.Debug("list tasks returned details", zap.Array("tasks", tasks))

	return tasks, nil
}

func (r *queryResolver) Login(ctx context.Context, id string, pwd string) (LoginRes, error) {
	user, err := r.Auth.UserFromContext(ctx, auth.OptAllowAll)
	if err != nil {
		return nil, err
	}

	r.Log.Debug("login...", zap.String("id", id))
	user, err = r.Dbal.GetUserByEmail(ctx, id)
	if err != nil {
		r.Log.Debug("")
		return Error{
			Code:    403,
			Message: "wrong username or/and password",
		}, err
	}

	token, err := r.Auth.GenerateToken(user.ID)
	if err != nil {
		r.Log.Error("cannot generate token", zap.String("user_id", user.ID), zap.Error(err))
		return Error{
			Code:    13,
			Message: "internal error",
		}, err
	}

	r.Log.Info("user successfully authenticated, token returned", zap.String("email", user.Email))

	return Token{
		Username: user.Email,
		Token:    token,
		UserID:   user.ID,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
