package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	wait := time.Millisecond * 1100
	t.Run("remove all user1 tasks", func(t *testing.T) {
		err := dbal.DeleteUserTasks(ctx, task1.UserId)
		assert.NoError(t, err)
	})

	<-time.After(wait)
	t.Run("create one", func(t *testing.T) {
		err := dbal.CreateTask(ctx, task1)
		assert.NoError(t, err)
	})

	id := ""
	<-time.After(wait)
	t.Run("get one by user id", func(t *testing.T) {
		res, err := dbal.GetTasksByUserID(ctx, task1.UserId)
		assert.Len(t, res, 1)
		assert.NoError(t, err)
		if len(res) > 0 {
			assertTask(t, task1, res[0])
			id = res[0].ID
		}
	})

	<-time.After(wait)
	if id != "" {
		t.Run("remove it", func(t *testing.T) {
			t.Log("try to remove ", id)
			err := dbal.DeleteTask(ctx, id)
			assert.NoError(t, err)
		})
	}

	<-time.After(wait)
	t.Run("get one by user id zero", func(t *testing.T) {
		res, err := dbal.GetTasksByUserID(ctx, task1.UserId)
		assert.Len(t, res, 0)
		assert.NoError(t, err)
	})
}

var task1 = Task{
	ID:        "id1",
	UserId:    "user1",
	CreatedAt: now,
	Status:    StatusNotStarted,
	Job: Job{
		DockerImage: "docker-img1",
		Dataset:     "data_set",
	},
}

var tasks = []Task{
	{
		ID:        "id3",
		UserId:    "user3",
		CreatedAt: now,
		Status:    StatusNotStarted,
		Job: Job{
			DockerImage: "docker-img3",
			Dataset:     "data_set3",
		},
	},
	{
		ID:        "id2",
		UserId:    "user2",
		CreatedAt: now,
		Status:    StatusNotStarted,
		Job: Job{
			DockerImage: "docker-img2",
			Dataset:     "data_set",
		},
	},
}

func assertTask(t *testing.T, expected, got Task) {
	if expected.UserId != got.UserId {
		t.Errorf("expected %v, got %v", expected.UserId, got.UserId)
	}
	if expected.ID != got.ID {
		t.Errorf("expected %v, got %v", expected.ID, got.ID)
	}
	if expected.Status != got.Status {
		t.Errorf("expected %v, got %v", expected.Status, got.Status)
	}
}
