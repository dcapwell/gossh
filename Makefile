all: clean build test

BASE = github.com/dcapwell/gossh
#PROTO_PKG = $(BASE)/prototype $(BASE)/workpool $(BASE)
PROTO_PKG = $(BASE)/workpool $(BASE)

clean:
	@go clean

build:
	@go build $(PROTO_PKG)

fmt:
	@go fmt $(PROTO_PKG)

test:
	@go test -v $(PROTO_PKG)
