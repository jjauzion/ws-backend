EXE = ws-backend

all: lint graphql build

graph/schema.resolvers.go: graph/schema.graphqls
	go run github.com/99designs/gqlgen generate

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: build
build:
	go build -o $(EXE)

.PHONY: run
run: all
	$(EXE) run

.PHONY: bootstrap
bootstrap: all
	$(EXE) run --bootstrap
