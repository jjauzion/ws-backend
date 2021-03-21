package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	wait := time.Millisecond * 1100
	id := ""
	t.Run("create user1", func(t *testing.T) {
		err := dbal.CreateUser(ctx, user1)
		assert.NoError(t, err)
	})

	<-time.After(wait)
	t.Run("get user1 by email", func(t *testing.T) {
		res, err := dbal.GetUserByEmail(ctx, user1.Email)
		assert.NoError(t, err)
		assertUser(t, user1, res)
		id = res.ID
	})

	<-time.After(wait)
	t.Run("get user1 by id", func(t *testing.T) {
		res, err := dbal.GetUserByID(ctx, id)
		assert.NoError(t, err)
		assertUser(t, user1, res)
	})

	<-time.After(wait)
	t.Run("delete user1", func(t *testing.T) {
		err := dbal.DeleteUser(ctx, id)
		assert.NoError(t, err)
	})

	<-time.After(wait)
	t.Run("get user1 by id error", func(t *testing.T) {
		_, err := dbal.GetUserByID(ctx, id)
		assert.Error(t, err)
	})
}

var user1 = User{
	ID:        "id-user1",
	Admin:     false,
	Email:     "user1@mail-adress.com",
	CreatedAt: now,
}

func assertUser(t *testing.T, expected, got User) {
	if expected.ID != got.ID {
		t.Errorf("expected %v, got %v", expected.ID, got.ID)
	}
	if expected.Email != got.Email {
		t.Errorf("expected %v, got %v", expected.Email, got.Email)
	}
}
