GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=$(shell basename `pwd`) 
BINARY_NAME_LINUX=$(shell basename `pwd`_linux) 

all:  test build

test:
	$(GOTEST) -v ./...
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
run:
	$(GOBUILD) -o $(BINARY_NAME) -v 
	./$(BINARY_NAME)
clean:
	$(GOLEAN)
	rm -f $(BINARY_NAME)

deps:
	$(GOGET) -v

#for MacOS or Windows
build-linux:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_LINUX) -v
