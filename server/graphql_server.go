package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jjauzion/ws-backend/db"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/graph"
)

func RunGraphQL(bootstrap bool) {
	ctx := context.Background()
	lg, cf, dbh, err := dependencies()
	if err != nil {
		return
	}

	resolver := &graph.Resolver{
		Log:     lg,
		DB:      dbh,
		Config:  cf,
		ApiPort: cf.WS_API_PORT,
	}

	if bootstrap {
		if err := db.Bootstrap(ctx, resolver.DB); err != nil {
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
