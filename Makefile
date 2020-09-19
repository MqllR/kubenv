BINARY := kubenv

.PONY: build
build:
		GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64
		GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64
