NAME := slp
GO_LDFLAGS=-s -w
GIT_COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build -trimpath -ldflags "$(GO_LDFLAGS) -X=main.version=$(GIT_COMMIT)" -o $(NAME) ./cmd/$(NAME)
