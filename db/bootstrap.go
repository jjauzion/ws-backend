package db

import (
	"github.com/google/uuid"
	"time"
)

func Bootstrap(dbh DatabaseHandler) error {
	err := dbh.Bootstrap()
	if err != nil {
		return err
	}
	userSimple := User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Email:     "simple-user@email.com",
		Admin:     false,
	}
	err = dbh.CreateUser(userSimple)
	if err != nil {
		return err
	}
	userAdmin := User{
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
	dockerImage := "ghcr.io/my-image"
	task1 := Task{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		StartedAt: time.Unix(0, 0),
		EndedAt:   time.Unix(0, 0),
		Status:    StatusNotStarted,
		UserId:    userAdmin.ID,
		Job:       Job{DockerImage: dockerImage, Dataset: dataset},
	}
	err = dbh.CreateTask(task1)
	if err != nil {
		return err
	}
	return err
}
