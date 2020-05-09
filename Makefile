PKGS=$(shell go list ./... | grep -v /vendor)
FMT_PKGS=$(shell go list -f {{.Dir}} ./... | grep -v vendor | grep -v test | tail -n +2)
PWD=$(shell pwd)
GOBUILD=go build -o ./bin/

default: authorsfile compile

all: default install

compile: ## Create tunnel executables in the ./bin directory
	${GOBUILD} ./cmd/tunnelctl
	${GOBUILD} ./cmd/tunneld

install: ## Create the tunnel executable in $GOPATH/bin directory.
	install -m 0755 bin/tunnelctl ${GOPATH}/bin/tunnelctl
	install -m 0755 bin/tunneld ${GOPATH}/bin/tunneld

clean: ## Clean the project tree from binary files.
	rm -rf bin/*

.PHONY: test
test: ## Run the tests.
	go test -v $(PKGS)

authorsfile: ## Update the AUTHORS file from the git logs
	git log -n 1 --format='%aN <%cE>' | sort -u | egrep -v "noreply|mailchimp|@Kris" > AUTHORS

vet: ## apply go vet to all the Go files
	@go vet $(PKGS)

gofmt: install-tools ## Go fmt your code
	echo "Fixing format of go files..."; \
	for package in $(FMT_PKGS); \
	do \
		gofmt -w $$package ; \
		goimports -l -w $$package ; \
	done

.PHONY: install-tools
install-tools:
	GOIMPORTS_CMD=$(shell command -v goimports 2> /dev/null)
ifndef GOIMPORTS_CMD
	go get golang.org/x/tools/cmd/goimports
endif

	GOLINT_CMD=$(shell command -v golint 2> /dev/null)
ifndef GOLINT_CMD
	go get golang.org/x/lint/golint
endif

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
