package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/db"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*User, error) {
	newUser := db.User{
		ID:        uuid.New().String(),
		Admin:     true,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}
	err := r.DB.CreateUser(ctx, newUser)
	if err != nil {
		if err == db.ErrTooManyRows {
			return nil, fmt.Errorf("user already exist")
		}
		r.Log.Warn("create user: ", zap.Error(err))
		return nil, err
	}

	return UserFromDBModel(newUser).Ptr(), nil
}

func (r *mutationResolver) CreateTask(ctx context.Context, input NewTask) (*Task, error) {
	mu, err := r.DB.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	user := UserFromDBModel(mu)
	newJob := db.Job{
		DockerImage: input.DockerImage,
		Dataset:     *input.Dataset,
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
	if err = r.DB.CreateTask(ctx, newTask); err != nil {
		return nil, err
	}
	r.Log.Info("task created", zap.String("id", newTask.ID))
	return TaskFromDBModel(newTask).Ptr(), nil
}

func (r *queryResolver) ListTasks(ctx context.Context, userID string) ([]*Task, error) {
	res, err := r.DB.GetTasksByUserID(ctx, userID)
	if err != nil {
		r.Log.Warn("cannot get tasks", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	tasks := Tasks{}
	for _, re := range res {
		tasks = append(tasks, TaskFromDBModel(re).Ptr())
	}

	r.Log.Debug("list tasks success", zap.Array("tasks", tasks))

	return tasks, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
