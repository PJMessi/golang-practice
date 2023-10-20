# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = myapp
GOTOOL = $(GOCMD) tool

all: build

build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -rf ./bin

run:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v 
	./bin/$(BINARY_NAME)

deps:
	$(GOGET) mod tidy

testunit:
	$(GOTEST) -v ./internal/... ./pkg/...

testintegration:
	$(GOTEST) -count=1 -v ./tests/...

test:
	$(GOTEST) -v ./...

testcov:
	$(GOTEST) -coverprofile=coverage.out ./...

testcovrep:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOTOOL) cover -html=coverage.out

.PHONY: all build clean run deps test
