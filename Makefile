EXE = ws-backend

all: lint graphql build
	echo "coucou"

.PHONY: graphql
graphql: graph/schema.graphqls
	go run github.com/99designs/gqlgen generate

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: build
build:
	go build -o $(EXE)
