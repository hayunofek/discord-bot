TEST?=$$(go list ./... | grep -v 'vendor')
NAME=discord-bot
BINARY=discord-bot
VERSION=0.0.1
OS_ARCH=linux_amd64

default: release

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64/${BINARY}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/${BINARY}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/${BINARY}.exe
