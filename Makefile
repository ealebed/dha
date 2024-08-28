VERSION=$(shell date -u '+%y.%m.%d-%H.%M')

GO       := go
GOBUILD  := CGO_ENABLED=0 $(GO) build 
GOTEST   := $(GO) test -gcflags='-l' -p 3

.PHONY: build
build:
	$(GOBUILD) -o bin/dha ./

.PHONY: install
install:
	$(GO) install ./

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: update
update:
	$(GO) get -u -v ./
	$(GO) mod verify
	$(GO) mod tidy
