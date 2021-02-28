package db

import (
	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/graph/model"
	"time"
)

func Bootstrap(dbh DatabaseHandler) error {
	err := dbh.Bootstrap()
	if err != nil {
		return err
	}
	userSimple := model.User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Email:     "simple-user@email.com",
		Admin:     false,
	}
	err = dbh.CreateUser(userSimple)
	if err != nil {
		return err
	}
	userAdmin := model.User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Email:     "admin-user@email.com",
		Admin:     true,
	}
	err = dbh.CreateUser(userAdmin)
	if err != nil {
		return err
	}
	dataset := "s3://task1"
	task1 := model.Task{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0, 0),
		EndedAt:   time.Unix(0, 0),
		Status:    model.StatusNotStarted,
		CreatedBy: userAdmin.ID,
		Job:       &model.Job{ID: uuid.New().String(), CreatedBy: userAdmin.ID, Dataset: &dataset},
	}
	err = dbh.CreateTask(task1)
	if err != nil {
		return err
	}
	return err
}
