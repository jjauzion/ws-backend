package server

import (
	"fmt"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/logger"
)

func buildDependencies() (logger.Logger, conf.Configuration, db.Dbal, error) {
	lg, err := logger.ProvideLogger()
	if err != nil {
		return logger.Logger{}, conf.Configuration{}, nil, fmt.Errorf("cannot create logger: %w", err)
	}

	cf, err := conf.GetConfig(lg)
	if err != nil {
		return logger.Logger{}, conf.Configuration{}, nil, fmt.Errorf("cannot get config: %w", err)
	}

	dbal, err := db.NewDatabaseAbstractedLayerImplemented(lg, cf)
	if err != nil {
		return logger.Logger{}, conf.Configuration{}, nil, fmt.Errorf("cannot create dbal: %w", err)
	}

	lg.Info("successfully connected to ES")

	return lg, cf, dbal, nil
}
