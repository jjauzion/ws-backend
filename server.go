package main

import (
	"fmt"
	"github.com/jjauzion/ws-backend/conf"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/graph"
	"github.com/jjauzion/ws-backend/graph/generated"
	"github.com/jjauzion/ws-backend/internal/logger"
)

const defaultPort = "8080"

func main() {

	resolver, err := Dependencies()
	if err != nil {
		return
	}

	//if err := resolver.DB.Bootstrap(); err != nil {
	//	return
	//}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	resolver.Log.Info(fmt.Sprintf("connect to http://localhost:%s/playground for GraphQL playground", resolver.ApiPort))
	resolver.Log.Error("", zap.Error(http.ListenAndServe(":"+resolver.ApiPort, nil)))
}

func Dependencies() (*graph.Resolver, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

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
		ApiPort: port,
	}

	return ret, nil
}
