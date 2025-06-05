build/linux:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./trader main.go

build/windows:
	GOARCH=amd64 GOOS=windows go build -trimpath -o ./trader.exe main.go

build: clean build/linux build/windows

clean:
	rm -f trader
	rm -f trader.exe
