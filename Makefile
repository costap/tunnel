PKGS=$(shell go list ./... | grep -v /vendor)
PWD=$(shell pwd)
GOBUILD=go build -o bin/

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

authorsfile: ## Update the AUTHORS file from the git logs
	git log --all --format='%aN <%cE>' | sort -u | egrep -v "noreply|mailchimp|@Kris" > AUTHORS

vet: ## apply go vet to all the Go files
	@go vet $(PKGS)