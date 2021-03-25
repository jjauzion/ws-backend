package server

import (
	"fmt"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/internal/auth"
	"github.com/jjauzion/ws-backend/internal/logger"
)

type application struct {
	log  logger.Logger
	conf conf.Configuration
	dbal db.Dbal
	auth auth.Auth
}

func buildApplication() (application, *graph.Resolver, error) {
	app := application{}
	var err error

	app.conf, err = conf.GetConfig()
	if err != nil {
		return app, nil, fmt.Errorf("cannot get config: %w", err)
	}

	app.log, err = logger.ProvideLogger(app.conf.Dev)
	if err != nil {
		return app, nil, fmt.Errorf("cannot create logger: %w", err)
	}

	app.dbal, err = db.NewDatabaseAbstractedLayerImplemented(app.log, app.conf)
	if err != nil {
		return app, nil, fmt.Errorf("cannot create dbal: %w", err)
	}

	app.auth, err = auth.NewAuth(app.dbal, app.log, app.conf.JWT_SIGNIN_KEY)
	if err != nil {
		return app, nil, fmt.Errorf("cannot initialize auth: %w", err)
	}

	return app, &graph.Resolver{
		Log:     app.log,
		Dbal:    app.dbal,
		Config:  app.conf,
		ApiPort: app.conf.WS_API_PORT,
		Auth:    app.auth,
	}, nil
}
