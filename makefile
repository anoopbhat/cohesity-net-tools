#Goparameters
GOCMD=/usr/local/go/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=coh-net-tools

all: test build
pkg:
	$(GOBUILD) -o $(BINARY_NAME) -v
	cp -a static $(GOPATH)/pkg/
	cp $(BINARY_NAME) $(GOPATH)/pkg/
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
