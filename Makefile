.PHONY: lint cover

lint:
	golangci-lint run ./...

cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic && \
	go tool cover -html=coverage.out -o coverage.html && \
	go tool cover -func=coverage.out | tail -1 && \
	rm -f coverage.out