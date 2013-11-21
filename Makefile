all: clean build test install

BASE = github.com/dcapwell/gossh
#PROTO_PKG = $(BASE)/prototype $(BASE)/workpool $(BASE)
PROTO_PKG = $(BASE)/workpool $(BASE)/utils $(BASE)

clean:
	@go clean

build:
	@go build $(PROTO_PKG)

install:
	@go install $(BASE)/gossh $(BASE)/gossh_agent
	@go install $(BASE)/gossh_agent

fmt:
	@go fmt $(PROTO_PKG)

test:
	@go test -v $(PROTO_PKG)

doc:
	@godoc -http=":6060" -index=true -play=true
