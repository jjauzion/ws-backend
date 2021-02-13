package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

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
	newUser := &model.User{
		ID:    input.UserID,
		Login: "toto",
		Email: "toto@a.com",
	}
	newTask := &model.Task{
		ID:        "1",
		CreatedBy: newUser,
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0, 0),
		EndedAt:    time.Unix(0, 0),
		Failed:    false,
		Job:       input.Job,
	}
	return newTask, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
