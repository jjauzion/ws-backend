package db

import (
	"go.uber.org/zap/zapcore"
)

type Tasks []*Task

func (t *Task) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", t.ID)
	encoder.AddString("user_id", t.UserId)
	return nil
}

func (t Tasks) MarshalLogArray(encoder zapcore.ArrayEncoder) error {
	for _, task := range t {
		err := encoder.AppendObject(task)
		if err != nil {
			return err
		}
	}
	return nil
}
