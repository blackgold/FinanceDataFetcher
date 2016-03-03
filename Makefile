GO_FILES = $(shell find . -type f -name '*.go')
GOPATH = $(shell pwd):$(shell pwd)/vendor
VENDORS_PATH = $(shell pwd)/vendor
all: setup build

build: $(GO_FILES)
	@GOPATH=$(GOPATH) go build -o bin/finance

clean:
	rm -rf bin/* pkg/*

setup:
	@GOPATH=$(VENDORS_PATH) go get github.com/mattn/go-sqlite3
	@GOPATH=$(VENDORS_PATH) go get gopkg.in/gorp.v1
