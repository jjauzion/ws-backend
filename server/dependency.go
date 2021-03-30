package server

import (
	"fmt"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/internal/auth"
	"github.com/jjauzion/ws-backend/internal/logger"
	"go.uber.org/zap"
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
		panic("cannot get config: " + err.Error())
	}

	app.log, err = logger.ProvideLogger(app.conf.ENV == "dev")
	if err != nil {
		panic("cannot create logger: " + err.Error())
	}

	app.dbal, err = db.NewDatabaseAbstractedLayerImplemented(app.log, app.conf)
	if err != nil {
		app.log.Error("cannot create dbal", zap.Error(err))
		return app, nil, fmt.Errorf("cannot create dbal: %w", err)
	}

	app.auth, err = auth.NewAuth(app.dbal, app.log, app.conf.JWT_SIGNIN_KEY, app.conf.TOKEN_DURATION_HOURS)
	if err != nil {
		app.log.Error("cannot initialize auth", zap.Error(err))
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
