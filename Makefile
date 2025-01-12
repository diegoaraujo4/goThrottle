all:
	go test ./...
	golangci-lint run
	go test ./... -coverprofile=coverage.txt