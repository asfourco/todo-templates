#!/usr/bin/make -f

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

build: go.sum
ifeq ($(OS),Windows_NT)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/todo ./backend/cmd/todo
endif

clean:
	rm -rf build/
install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./backend/cmd/todo

lint:
	golangci-lint run --tests=false
format:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -w -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -w -local github.com/