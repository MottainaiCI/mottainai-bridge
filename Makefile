NAME ?= mottainai-bridge
PACKAGE_NAME ?= $(NAME)
PACKAGE_CONFLICT ?= $(PACKAGE_NAME)-beta
REVISION := $(shell git rev-parse --short HEAD || echo unknown)
VERSION := $(shell git describe --tags || cat pkg/settings/settings.go | echo $(REVISION) || echo dev)
VERSION := $(shell echo $(VERSION) | sed -e 's/^v//g')
ITTERATION := $(shell date +%s)
BUILD_PLATFORMS ?= -osarch="linux/amd64" -osarch="linux/386" -osarch="linux/arm"
SUBDIRS =
DESTDIR =
UBINDIR ?= /usr/bin
LIBDIR ?= /usr/lib
SBINDIR ?= /sbin
USBINDIR ?= /usr/sbin
BINDIR ?= /bin
LIBEXECDIR ?= /usr/libexec
SYSCONFDIR ?= /etc
LOCKDIR ?= /var/lock
LIBDIR ?= /var/lib
EXAMPLES = simplelistener

all: deps multiarch-build install

build-test: test build

help:
	# make all => deps test lint build
	# make deps - install all dependencies
	# make test - run project tests
	# make lint - check project code style
	# make build - build project for all supported OSes

clean:
	rm -rf vendor/
	rm -rf release/

deps:
	go env
	# Installing dependencies...
	go get -u github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

build:
	for d in $(EXAMPLES); do pushd examples/$$d; go build; popd; done

lint:
	# Checking project code style...
	golint ./... | grep -v "be unexported"

test:
	# Running tests... ${TOTEST}
	go test -v -cover -race ./...
