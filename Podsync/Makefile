BINPATH := $(abspath ./bin)
GOLANGCI := $(BINPATH)/golangci-lint

.PHONY: all
all: build lint test

#
# Build Podsync CLI binary
#
.PHONY: build
build:
	go build -o podsync ./cmd/podsync

#
# Build Docker image
#
.PHONY: docker
docker:
	GOOS=linux GOARCH=amd64 go build -o podsync ./cmd/podsync
	docker build -t mxpv/podsync:unstable .
	docker push mxpv/podsync:unstable

#
# Run goreleaser to build and upload release binaries
#
.PHONY: release
release:
	goreleaser --rm-dist

#
# Pull GolangCI-Lint dependency
#
$(GOLANGCI):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BINPATH) v1.17.1
	$(GOLANGCI) --version

#
# Run linter
#
.PHONY: lint
lint: $(GOLANGCI)
	$(GOLANGCI) run

#
# Run unit tests
#
.PHONY: test
test:
	go test -v ./...
