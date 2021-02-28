package server

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/internal/logger"
)

func Run(bootstrap bool) {
	resolver, err := Dependencies()
	if err != nil {
		return
	}

	if bootstrap {
		if err := db.Bootstrap(resolver.DB); err != nil {
			resolver.Log.Error("bootstrap failed", zap.Error(err))
			return
		}
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	resolver.Log.Info(fmt.Sprintf("connect to http://localhost:%s/playground for GraphQL playground", resolver.ApiPort))
	resolver.Log.Error("", zap.Error(http.ListenAndServe(":"+resolver.ApiPort, nil)))
}

func Dependencies() (*graph.Resolver, error) {
	lg, err := logger.ProvideLogger()
	if err != nil {
		log.Fatalf("cannot create logger %v", err)
	}

	cf, err := conf.GetConfig(lg)
	if err != nil {
		log.Fatalf("cannot get config %v", err)
	}

	dbh := db.NewDBHandler(lg, cf)
	ret := &graph.Resolver{
		Log:     lg,
		DB:      dbh,
		Config:  cf,
		ApiPort: cf.WS_API_PORT,
	}

	return ret, nil
}
