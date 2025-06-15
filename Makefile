build/linux:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./trader main.go

build/windows:
	GOARCH=amd64 GOOS=windows go build -trimpath -o ./trader.exe main.go

build/darwin:
	GOARCH=arm64 GOOS=darwin go build -trimpath -o ./trader main.go

test/unit:
	go test -coverprofile=coverage.out ./...

test/coverage:
	go tool cover -html=coverage.out

test/all: test/unit test/coverage

