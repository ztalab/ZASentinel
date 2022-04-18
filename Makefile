.PHONY: start build

APP 			= zero-access
SERVER_BIN  	= ./bin/backend
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(GIT_COUNT).$(GIT_HASH)

GOPROXY=https://goproxy.oneitfarm.com,https://goproxy.cn,direct

CFLAGS = -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) ./cmd

all: start

build:
	@go build $(CFLAGS)

start:
	go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./cmd/main.go web -c ./configs/config.toml

test:
	@go test -v $(shell go list ./...)

clean:
	rm -rf data release $(SERVER_BIN) internal/app/test/data cmd/data

release:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(CFLAGS)