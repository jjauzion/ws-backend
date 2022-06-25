package graph

import (
	"github.com/42-AI/ws-backend/db"
)

func TaskFromDBModel(m db.Task) Task {
	return Task{
		ID:        m.ID,
		UserID:    m.UserId,
		CreatedAt: m.CreatedAt,
		StartedAt: m.StartedAt,
		EndedAt:   m.EndedAt,
		Status:    Status(m.Status.String()),
		Job:       JobFromDBModel(m.Job),
	}
}

func JobFromDBModel(j db.Job) *Job {
	return &Job{
		DockerImage: j.DockerImage,
		Dataset:     &j.Dataset,
		Env:         j.Env,
	}
}

func UserFromDBModel(u db.User) User {
	return User{
		ID:        u.ID,
		Admin:     u.Admin,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func TaskToDBModel(m Task) db.Task {
	return db.Task{
		ID:        m.ID,
		UserId:    m.UserID,
		CreatedAt: m.CreatedAt,
		StartedAt: m.StartedAt,
		EndedAt:   m.EndedAt,
		Status:    db.Status(m.Status.String()),
		Job:       JobToDBModel(*m.Job),
	}
}

func JobToDBModel(j Job) db.Job {
	return db.Job{
		DockerImage: j.DockerImage,
		Dataset:     *j.Dataset,
		Env:         j.Env,
	}
}

func UserToDBModel(u User) db.User {
	return db.User{
		ID:        u.ID,
		Admin:     u.Admin,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func (j Job) Ptr() *Job {
	return &j
}

func (t Task) Ptr() *Task {
	return &t
}

func (u User) Ptr() *User {
	return &u
}
