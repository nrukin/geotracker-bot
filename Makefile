BINARY_NAME=geotracker-bot

build:
	go build -o ${BINARY_NAME} .
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux-amd .
	GOARCH=arm64 GOOS=linux go build -o ${BINARY_NAME}-linux-arm .

run: build
	./${BINARY_NAME}
