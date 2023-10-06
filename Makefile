# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = myapp

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

test:
	$(GOTEST) -v ./...

.PHONY: all build clean run deps test
