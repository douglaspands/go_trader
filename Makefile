build/linux:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./trader main.go

build/windows:
	GOARCH=amd64 GOOS=windows go build -trimpath -o ./trader.exe main.go

build/darwin:
	GOARCH=arm64 GOOS=darwin go build -trimpath -o ./trader main.go

deps/tidy:
	go mod tidy

deps/sync: deps/tidy
deps/cleanup: deps/tidy

test/unit:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test/coverage:
	go tool cover -html=coverage.out

test/unit/coverage: test/unit test/coverage
