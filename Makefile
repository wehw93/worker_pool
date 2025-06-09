.PHONY: build
build:
	go build -v ./cmd/worker_pool/main.go

.DEFAULT_GOAL := build
