CURDIR = $(shell pwd)
GOPATH = "$(CURDIR)/.gopath"
all: build

# `make setup` to set up git submodules
setup:
	git submodule init
	git submodule update

# `make run` to build and run the bot
run: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go run cmd/scarecrow/main.go

# `make debug` to build and run in debug mode
debug: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go run cmd/scarecrow/main.go --debug

# `make fmt` runs gofmt
fmt:
	gofmt -w .

# `make build` to build the binary
build:
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 \
		go build -x -o bin/scarecrow cmd/scarecrow/main.go

# Sets up the gopath / build environment
gopath:
	export GO15VENDOREXPERIMENT=1
	mkdir -p .gopath/src/github.com/aichaos/scarecrow bin
	ln -sf "$(CURDIR)/src" .gopath/src/github.com/aichaos/scarecrow
	ln -sf "$(CURDIR)/vendor" .gopath/src/github.com/aichaos/scarecrow/src
