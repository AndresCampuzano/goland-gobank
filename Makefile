build:
	@go build -o bin/goland-gobank

run: build
	@./bin/goland-gobank

test:
	@go test -v ./...