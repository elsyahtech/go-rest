.PHONY: run lint test coverage check clean build ci

lint:
	@echo "Running only golangci-lint..."
	golangci-lint run ./...

upgrade:
	@echo "Running uprade dependency..."
	go get -u ./...
	go mod tidy