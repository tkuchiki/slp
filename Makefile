NAME := slp
GO_LDFLAGS=-s -w

.PHONY: build
build:
	go build -trimpath -ldflags "$(GO_LDFLAGS)" -o $(NAME) ./cmd/$(NAME)
