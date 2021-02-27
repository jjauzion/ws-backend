package db

import (
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScene1(t *testing.T) {
	wait := time.Millisecond * 900
	t.Run("remove all user1 tasks", func(t *testing.T) {
		err := dbh.DeleteUserTasks(ctx, task1.CreatedBy)
		assert.NoError(t, err)
	})

	<-time.After(wait)
	t.Run("create one", func(t *testing.T) {
		err := dbh.CreateTask(task1)
		assert.NoError(t, err)
	})

	id := ""
	<-time.After(wait)
	t.Run("get one by user id", func(t *testing.T) {
		res, err := dbh.GetTasksByUserID(ctx, task1.CreatedBy)
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
			err := dbh.DeleteTask(ctx, id)
			assert.NoError(t, err)
		})
	}

	<-time.After(wait)
	t.Run("get one by user id zero", func(t *testing.T) {
		res, err := dbh.GetTasksByUserID(ctx, task1.CreatedBy)
		assert.Len(t, res, 0)
		assert.NoError(t, err)
	})
}

var task1 = model.Task{
	ID:        "id1",
	CreatedBy: "user1",
	CreatedAt: now,
	Failed:    false,
	Job: &model.Job{
		ID:          "id1",
		CreatedBy:   "user1",
		DockerImage: "docker-img1",
		Dataset:     toStringPtr("data_set"),
	},
}

var tasks = []model.Task{
	{
		ID:        "id3",
		CreatedBy: "user3",
		CreatedAt: now,
		Failed:    false,
		Job: &model.Job{
			ID:          "id3",
			CreatedBy:   "user3",
			DockerImage: "docker-img3",
			Dataset:     toStringPtr("data_set3"),
		},
	},
	{
		ID:        "id2",
		CreatedBy: "user2",
		CreatedAt: now,
		Failed:    false,
		Job: &model.Job{
			ID:          "id2",
			CreatedBy:   "user2",
			DockerImage: "docker-img2",
			Dataset:     toStringPtr("data_set"),
		},
	},
}

func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func assertTask(t *testing.T, expected, got model.Task) {
	if expected.CreatedBy != got.CreatedBy {
		t.Errorf("expected %v, got %v", expected.CreatedBy, got.CreatedBy)
	}
	if expected.CreatedAt != got.CreatedAt {
		t.Errorf("expected %v, got %v", expected.CreatedAt, got.CreatedAt)
	}
	if expected.ID != got.ID {
		t.Errorf("expected %v, got %v", expected.ID, got.ID)
	}
	if expected.Failed != got.Failed {
		t.Errorf("expected %v, got %v", expected.Failed, got.Failed)
	}
}
