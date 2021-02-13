package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/pkg"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log := pkg.NewLog()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", port))
	log.Error("", zap.Error(http.ListenAndServe(":"+port, nil)))
}
