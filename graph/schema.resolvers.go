package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/graph/model"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	newUser := model.User{
		ID:        uuid.New().String(),
		Admin:     true,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}
	if err := r.DB.CreateUser(newUser); err != nil {
		if err == db.ErrTooManyRows {
			return nil, fmt.Errorf("user already exist")
		}
		r.Log.Warn("create user: ", zap.Error(err))
		return nil, err
	}

	return &newUser, nil
}

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	var user model.User
	var err error
	if user, err = r.DB.GetUserByID(input.UserID); err != nil {
		return nil, err
	}
	newJob := &model.Job{
		ID:          uuid.New().String(),
		CreatedBy:   user.ID,
		DockerImage: input.DockerImage,
		Dataset:     input.Dataset,
	}
	newTask := &model.Task{
		ID:        uuid.New().String(),
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0, 0),
		EndedAt:   time.Unix(0, 0),
		Failed:    false,
		Job:       newJob,
	}
	if err = r.DB.CreateTask(*newTask); err != nil {
		return nil, err
	}
	r.Log.Info("task created", zap.String("id", newTask.ID))
	return newTask, nil
}

func (r *queryResolver) ListTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	res, err := r.DB.GetTasksByUserID(ctx, userID)
	if err != nil {
		r.Log.Warn("cannot get tasks", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	tasks := model.Tasks{}
	for _, re := range res {
		tasks = append(tasks, &re)
	}

	r.Log.Debug("list tasks success", zap.Array("tasks", tasks))

	return tasks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
