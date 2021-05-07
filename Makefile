# The crossbuild command sets the PREFIX env var for each platform
PREFIX ?= $(shell pwd)

GO           ?= go
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
PROMU        := $(FIRST_GOPATH)/bin/promu

.PHONY: build
build: tools.promu
	$(PROMU) build --prefix $(PREFIX)

.PHONY: crossbuild
crossbuild:
	$(PROMU) crossbuild
	$(PROMU) crossbuild tarballs

.PHONY: release
release: crossbuild
	# this should be invoked after drafting the release in github.com
	$(PROMU) release .tarballs

#---------------
#-- tools
#---------------
.PHONY: tools
tools: tools.golangci-lint tools.promu

.PHONY: tools.golangci-lint
tools.golangci-lint:
	@command -v golangci-lint >/dev/null || { \
		echo ">> Installing golangci-lint..."; \
		GOOS= GOARCH= $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint; \
	}

.PHONY: tools.promu
tools.promu:
	@command -v promu >/dev/null || { \
		echo ">> installing promu"; \
		GOOS= GOARCH= $(GO) install github.com/prometheus/promu; \
	}
