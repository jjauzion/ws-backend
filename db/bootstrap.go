package db

import (
	"context"
	"github.com/google/uuid"
	"time"
)

func Bootstrap(ctx context.Context, dbal Dbal) error {
	err := dbal.CreateIndexes(ctx)
	if err != nil {
		return err
	}
	userSimple := User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Email:     "simple-user@email.com",
		Admin:     false,
	}
	err = dbal.CreateUser(ctx, userSimple)
	if err != nil {
		return err
	}
	userAdmin := User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Email:     "admin-user@email.com",
		Admin:     true,
	}
	err = dbal.CreateUser(ctx, userAdmin)
	if err != nil {
		return err
	}
	return err
}
