BINARY := kubenv
LDFLAGS := -X 'github.com/mqllr/kubenv/cmd.Version=$(VERSION)'

all: test build

build:
		GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-linux-amd64
		GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-darwin-amd64

test: install_deps
	go test -v ./...

lint:
	golangci-lint run

install_deps:
	go get -v ./...

clean:
	rm -f $(BINARY)-*-amd64
