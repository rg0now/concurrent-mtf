BINARY_NAME=cmtf

build: dep
	GOARCH=amd64 GOOS=linux go build -tags netgo -o ${BINARY_NAME} *.go

dep:
	go mod tidy -go=1.18 -compat=1.18

all: build

clean:
	go clean
	rm ${BINARY_NAME}
