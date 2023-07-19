build: $(wildcard *.go */*.go */*/*.go)
	@echo Building
	@go build -o busyboi .

run: $(wildcard *.go */*.go */*/*.go)
	@echo Running
	@go run .

test:
	@echo Testing
	@go test -short ./...
