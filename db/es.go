package db

import (
	"context"
	"github.com/jjauzion/ws-backend/pkg"
	"go.uber.org/zap"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)


func NewClient(ctx context.Context) (es *elasticsearch7.Client, err error) {
	log := pkg.GetLogger()
	if es, err = elasticsearch7.NewDefaultClient(); err != nil {
		return
	}
		log.Error("couldn't connect to DB:", zap.Error(err))
	return

}