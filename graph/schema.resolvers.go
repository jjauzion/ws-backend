package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/jjauzion/ws-backend/db"
	"time"

	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/jjauzion/ws-backend/internal"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	dbh := db.NewDBHandler()
	newUser := model.User{
		ID:    uuid.New().String(),
		Admin: true,
		Email: input.Email,
	}
	if err := dbh.CreateUser(newUser); err != nil {
		return nil, err
	}
	return &newUser, nil
}

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	log := internal.GetLogger()
	newUser := &model.User{
		ID:    input.UserID,
		Admin: true,
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
