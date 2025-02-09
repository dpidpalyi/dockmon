DEFAULT_GOAL: build

.PHONY: fmt vet build clean run imports

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build ./cmd/web

run: vet
	go run ./cmd/web

imports:
	find . -name \*.go -exec goimports -w -l {} \;

clean:
	rm -rf web
