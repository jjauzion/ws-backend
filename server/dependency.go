package server

import (
	"fmt"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/auth"
	"github.com/jjauzion/ws-backend/internal/logger"
)

type application struct {
	log  logger.Logger
	conf conf.Configuration
	dbal db.Dbal
	auth auth.Auth
}

func buildApplication() (application, error) {
	app := application{}
	var err error

	app.log, err = logger.ProvideLogger()
	if err != nil {
		return app, fmt.Errorf("cannot create logger: %w", err)
	}

	app.conf, err = conf.GetConfig(app.log)
	if err != nil {
		return app, fmt.Errorf("cannot get config: %w", err)
	}

	app.dbal, err = db.NewDatabaseAbstractedLayerImplemented(app.log, app.conf)
	if err != nil {
		return app, fmt.Errorf("cannot create dbal: %w", err)
	}

	app.auth, err = auth.NewAuth(app.dbal, app.log, app.conf.JWT_SIGNIN_KEY)
	if err != nil {
		return app, fmt.Errorf("cannot initialize auth: %w", err)
	}

	return app, nil
}
