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
	app, resolver, err := buildApplication()
	if err != nil {
		log.Fatal(err)
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

	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	resolver.Log.Info(fmt.Sprintf("connect to http://localhost:%s/playground for GraphQL playground", resolver.ApiPort))
	resolver.Log.Error("", zap.Error(http.ListenAndServe(":"+resolver.ApiPort, router)))
}
