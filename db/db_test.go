package db

import (
	"context"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

var dbal Dbal
var ctx = context.Background()
var now = time.Now()

func TestMain(m *testing.M) {
	address := "http://localhost:9200"
	lg := zap.NewNop()

	elst, err := elastic.NewClient(elastic.SetURL(address),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		panic(err)
	}

	dbal = &esHandler{
		conf:   conf.Configuration{},
		log:    logger.Logger{Logger: lg},
		client: elst,
	}

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}

func TestEsHandler_Info(t *testing.T) {
	err := dbal.Ping()
	if err != nil {
		t.Error(err)
	}
}
