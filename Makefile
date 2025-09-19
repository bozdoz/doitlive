build-all: build build-linux

# for smaller binaries: 
# CGO_ENABLED=0 go build -ldflags="-s -w" .
build:
	go build -o dist/doitlive .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/doitlive-linux .