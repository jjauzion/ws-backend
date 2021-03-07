package server

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jjauzion/ws-backend/db"
	"go.uber.org/zap"
	"net/http"

	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/graph/generated"
)

func RunGraphQL(bootstrap bool) {
	lg, cf, dbh, err := dependencies()
	resolver := &graph.Resolver{
		Log:     lg,
		DB:      dbh,
		Config:  cf,
		ApiPort: cf.WS_API_PORT,
	}
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
