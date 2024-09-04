build:
	@go build -o bin/blockr

run: build
	@./bin/blockr

test:
	@go test -v ./...
