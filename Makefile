BINARY_NAME=cmtf

build:
	GOARCH=amd64 GOOS=linux go build -tags netgo -o ${BINARY_NAME} *.go

all: build

clean:
	go clean
	rm ${BINARY_NAME}
