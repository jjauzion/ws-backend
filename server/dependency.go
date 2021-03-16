package server

import (
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/olivere/elastic/v7"
	"log"
)

func buildDependencies() (*logger.Logger, conf.Configuration, db.Dbal, error) {
	lg, err := logger.ProvideLogger()
	if err != nil {
		log.Fatalf("cannot create logger %v", err)
	}

	cf, err := conf.GetConfig(lg)
	if err != nil {
		log.Fatalf("cannot get config %v", err)
	}

	elasticClient, err := elastic.NewClient(elastic.SetURL(cf.WS_ES_HOST+":"+cf.WS_ES_PORT),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		log.Fatalf("cannot create elastic client %v", err)
	}

	dbal := db.NewDatabaseAbstractedLayerImplemented(lg, cf, elasticClient)

	return lg, cf, dbal, nil
}
