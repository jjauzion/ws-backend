package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"

	"github.com/jjauzion/ws-backend/pkg"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	newUser := &model.User{
		ID:    "1",
		Login: input.Login,
		Email: input.Email,
	}
	return newUser, nil
}

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	log := pkg.NewLog()
	newUser := &model.User{
		ID:    input.UserID,
		Login: "toto",
		Email: "toto@a.com",
	}
	newJob := &model.Job{
		ID:          uuid.New().String(),
		CreatedBy:   newUser,
		DockerImage: input.DockerImage,
		Dataset:     input.Dataset,
	}
	newTask := &model.Task{
		ID:        "1",
		CreatedBy: newUser,
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0, 0),
		EndedAt:   time.Unix(0, 0),
		Failed:    false,
		Job:       newJob,
	}
	log.Info("task created", zap.String("id", newTask.ID))
	return newTask, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
