package server

import (
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/olivere/elastic/v7"
	"log"
)

func dependencies() (*logger.Logger, conf.Configuration, db.DatabaseHandler, error) {
	lg, err := logger.ProvideLogger()
	if err != nil {
		log.Fatalf("cannot create logger %v", err)
	}

	cf, err := conf.GetConfig(lg)
	if err != nil {
		log.Fatalf("cannot get config %v", err)
	}

	elst, err := elastic.NewClient(elastic.SetURL(cf.WS_ES_HOST+":"+cf.WS_ES_PORT),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	dbh := db.NewDBHandler(lg, cf, elst)

	return lg, cf, dbh, nil
}
