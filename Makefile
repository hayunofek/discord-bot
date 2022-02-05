TEST?=$$(go list ./... | grep -v 'vendor')
NAME=discord-bot
BINARY=discord-bot
OS_ARCH=linux_amd64
VERSION=0.0.5

default: release

build:
	go build -o ${BINARY}

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/${BINARY}
