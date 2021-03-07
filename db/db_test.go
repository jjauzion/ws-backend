package db

import (
	"context"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

var dbh DatabaseHandler
var ctx = context.Background()
var now = time.Now()

func TestMain(m *testing.M) {
	address := "http://localhost:9200"
	esConfig := elasticsearch7.Config{
		Addresses: []string{address},
	}
	client, err := elasticsearch7.NewClient(esConfig)
	if err != nil {
		panic(err)
	}
	lg := zap.NewNop()

	elst, err := elastic.NewClient(elastic.SetURL(address),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	dbh = &esHandler{
		client:  client,
		cf:      conf.Configuration{},
		log:     &logger.Logger{Logger: lg},
		elastic: elst,
	}
	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}
