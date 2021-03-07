package server

import (
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/logger"
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

	dbh := db.NewDBHandler(lg, cf)

	return lg, cf, dbh, nil
}
