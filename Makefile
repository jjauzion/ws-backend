EXE = ./ws-backend
GRAPHQL_FILES = graph/schema.resolvers.go
SRC_FILES = $(wildcard *.go) \
            $(wildcard */*.go)

all: lint $(GRAPHQL_FILES) $(EXE)

$(EXE): $(SRC_FILES)
	go build -o $(EXE)

$(GRAPHQL_FILES): graph/schema.graphqls
	go run github.com/99designs/gqlgen generate

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: run
run: all
	$(EXE) run

.PHONY: bootstrap
bootstrap: all
	$(EXE) run --bootstrap
