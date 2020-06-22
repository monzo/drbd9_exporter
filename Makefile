PKG    = github.com/benjojo/drbd9_exporter
PREFIX = /usr

all: build/drbd9_exporter

# NOTE: This repo uses Go modules, and uses a synthetic GOPATH at
# $(CURDIR)/.gopath that is only used for the build cache. $GOPATH/src/ is
# empty.
GO            = GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS =
GO_LDFLAGS    = -s -w

APP_VERSION ?= v0.0.0-dev
GIT_REVISION ?= $(shell git rev-parse --short HEAD)
GIT_BRANCH ?= $(shell git symbolic-ref -q --short HEAD)
GIT_EMAIL ?= $(shell git config --get user.email)

GO_LDFLAGS += -X 'github.com/prometheus/common/version.Version=$(APP_VERSION)'
GO_LDFLAGS += -X 'github.com/prometheus/common/version.Revision=$(GIT_REVISION)'
GO_LDFLAGS += -X 'github.com/prometheus/common/version.Branch=$(GIT_BRANCH)'
GO_LDFLAGS += -X 'github.com/prometheus/common/version.BuildUser=$(GIT_EMAIL)'
GO_LDFLAGS += -X 'github.com/prometheus/common/version.BuildDate=$(shell date -u "+%Y-%m-%dT%H:%M:%S%z")'

build/drbd9_exporter: *.go
	$(GO) install $(GO_BUILDFLAGS) -ldflags "$(GO_LDFLAGS)" .

install: build/drbd9_exporter
	install -D -m 0755 build/drbd9_exporter "$(DESTDIR)$(PREFIX)/bin/drbd9_exporter"

vendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: install vendor