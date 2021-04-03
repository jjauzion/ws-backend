package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/graph/playground"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/graph"
)

func RunGraphQL(bootstrap bool) {
	ctx := context.Background()
	app, resolver, err := buildApplication()
	if err != nil {
		return
	}

	if bootstrap {
		err = db.Bootstrap(ctx, app.dbal)
		if err != nil {
			resolver.Log.Error("bootstrap failed", zap.Error(err))
			return
		}
	}

	router := chi.NewRouter()
	router.Use(app.auth.Middleware())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	router.Handle("/query", srv)
	if app.conf.IS_DEV_ENV {
		router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	}

	resolver.Log.Info(fmt.Sprintf("connect to http://localhost:%s/playground for GraphQL playground", resolver.ApiPort))
	resolver.Log.Error("", zap.Error(http.ListenAndServe(":"+resolver.ApiPort, router)))
}
