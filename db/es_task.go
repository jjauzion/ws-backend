package db

import (
	"bytes"
	"encoding/json"
	"github.com/jjauzion/ws-backend/graph/model"
	"go.uber.org/zap"
)

func (es *esHandler) CreateTask(task model.Task) (err error) {
	es.log.Debug("creating new task...")
	var b []byte
	if b, err = json.Marshal(task); err != nil {
		return
	}
	if err = es.indexNewDoc(taskIndex, bytes.NewReader(b)); err != nil {
		es.log.Error("failed to create task", zap.Error(err))
		return
	}
	es.log.Info("new task successfully created", zap.String("id", task.ID))
	return
}
