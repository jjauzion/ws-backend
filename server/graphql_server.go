package server

import (
	"context"
	"fmt"
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
	lg, conf, dbal, err := buildDependencies()
	if err != nil {
		log.Fatal(err)
	}

	resolver := &graph.Resolver{
		Log:     lg,
		Dbal:    dbal,
		Config:  conf,
		ApiPort: conf.WS_API_PORT,
	}

	if bootstrap {
		err := db.Bootstrap(ctx, dbal)
		if err != nil {
			resolver.Log.Error("bootstrap failed", zap.Error(err))
			return
		}
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	resolver.Log.Info(fmt.Sprintf("connect to http://localhost:%s/playground for GraphQL playground", resolver.ApiPort))
	resolver.Log.Error("", zap.Error(http.ListenAndServe(":"+resolver.ApiPort, nil)))
}
