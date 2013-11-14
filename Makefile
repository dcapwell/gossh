all: clean build test

PROTO_PKG = github.com/dcapwell/gossh/prototype/unix_socket

clean:
	@go clean

build:
	@go build $(PROTO_PKG)

fmt:
	@go fmt $(PROTO_PKG)

test:
	@go test -v $(PROTO_PKG)
