BINARY := kubenv

all: test build

build:
		GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64
		GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64

test: install_deps lint
	go test -v ./...

lint:
	golangci-lint run

install_deps:
	go get -v ./...

clean:
	rm -f $(BINARY)-*-amd64
