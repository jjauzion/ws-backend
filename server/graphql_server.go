package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jjauzion/ws-backend/db"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/graph"
)

func RunGraphQL(bootstrap bool) {
	ctx := context.Background()
	app, err := buildApplication()
	if err != nil {
		log.Fatal(err)
	}

	resolver := &graph.Resolver{
		Log:     app.log,
		Dbal:    app.dbal,
		Config:  app.conf,
		ApiPort: app.conf.WS_API_PORT,
	}

	if bootstrap {
		err := db.Bootstrap(ctx, app.dbal)
		if err != nil {
			resolver.Log.Error("bootstrap failed", zap.Error(err))
			return
		}
	}

	router := chi.NewRouter()

	router.Use(app.auth.Middleware())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	resolver.Log.Info(fmt.Sprintf("connect to http://localhost:%s/playground for GraphQL playground", resolver.ApiPort))
	resolver.Log.Error("", zap.Error(http.ListenAndServe(":"+resolver.ApiPort, nil)))
}
