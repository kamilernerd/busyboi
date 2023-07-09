build: $(wildcard *.go */*.go */*/*.go)
	@echo 🤸 go build !
	@go build -o busyboi .

run: $(wildcard *.go */*.go */*/*.go)
	@echo RUN!
	@go run .

fmtcheck:
	@echo 🦉 Checking format with gofmt -d -s
	@if [ "x$$(find . -name '*.go' -not -wholename './gen/*' -and -not -wholename './vendor/*' -exec gofmt -d -s {} +)" != "x" ]; then find . -name '*.go' -not -wholename './gen/*' -and -not -wholename './vendor/*' -exec gofmt -d -s {} +; exit 1; fi

fmtfix:
	@echo 🎨 Fixing formating
	@find . -name '*.go' -not -wholename './gen/*' -and -not -wholename './vendor/*' -exec gofmt -d -s -w {} +

test:
	@echo 🧐 Testing
	@go test -short ./...
