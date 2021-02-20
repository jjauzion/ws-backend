package main

import (
	"fmt"
	"github.com/jjauzion/ws-backend/db"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/internal"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log := internal.GetLogger()

	dbh := db.NewDBHandler()
	if err := dbh.Bootstrap(); err != nil {
		return
	}
	dbh.GetUserByEmail("test")

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))


	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", port))
	log.Error("", zap.Error(http.ListenAndServe(":"+port, nil)))
}
