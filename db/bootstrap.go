package db

import (
	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/jjauzion/ws-backend/internal/logger"
	"go.uber.org/zap"
	"time"
)

func Bootstrap(dbh DatabaseHandler, log *logger.Logger) error {
	err := dbh.Bootstrap()
	if err != nil {
		log.Error("bootstrap failed", zap.Error(err))
		return err
	}
	userSimple := model.User{
		ID: uuid.New().String(),
		CreatedAt: time.Now(),
		Email: "simple-user@email.com",
		Admin: false,
	}
	err = dbh.CreateUser(userSimple)
	if err != nil {
		log.Error("bootstrap failed", zap.Error(err))
		return err
	}
	userAdmin := model.User{
		ID: uuid.New().String(),
		CreatedAt: time.Now(),
		Email: "admin-user@email.com",
		Admin: true,
	}
	err = dbh.CreateUser(userAdmin)
	if err != nil {
		log.Error("bootstrap failed", zap.Error(err))
		return err
	}
	dataset := "s3://task1"
	task1 := model.Task{
		ID: uuid.New().String(),
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0,0),
		EndedAt: time.Unix(0,0),
		Status: model.StatusNotStarted,
		CreatedBy: userAdmin.ID,
		Job: &model.Job{ID: uuid.New().String(), CreatedBy: userAdmin.ID, Dataset: &dataset},
	}
	err = dbh.CreateTask(task1)
	if err != nil {
		log.Error("bootstrap failed", zap.Error(err))
		return err
	}
	return err
}
