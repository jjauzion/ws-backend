package db

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Job struct {
	DockerImage string `json:"docker_image"`
	Dataset     string `json:"dataset"`
}

func (j Job) Ptr() *Job {
	return &j
}

type Task struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	Status    Status    `json:"status"`
	Job       Job       `json:"job"`
}

func (t Task) Ptr() *Task {
	return &t
}

type User struct {
	ID        string    `json:"id"`
	Admin     bool      `json:"admin"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) Ptr() *User {
	return &u
}

type Status string

const (
	StatusStatusless Status = "STATUSLESS"
	StatusFailed     Status = "FAILED"
	StatusNotStarted Status = "NOT_STARTED"
	StatusRunning    Status = "RUNNING"
	StatusEnded      Status = "ENDED"
	StatusCanceled   Status = "CANCELED"
)

var AllStatus = []Status{
	StatusStatusless,
	StatusFailed,
	StatusNotStarted,
	StatusRunning,
	StatusEnded,
	StatusCanceled,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusStatusless, StatusFailed, StatusNotStarted, StatusRunning, StatusEnded, StatusCanceled:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
