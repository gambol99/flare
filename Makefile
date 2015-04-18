#
#   Author: Rohith
#   Date: 2015-04-18 11:54:16 +0100 (Sat, 18 Apr 2015)
#
#  vim:ts=2:sw=2:et
#
NAME="flare"
AUTHOR=gambol99
HARDWARE=$(shell uname -m)

default: build

.PHONY: clean test build changelog

clean:
	test -d ./bin && rm -rf ./bin || true

build: clean
	mkdir -p ./bin
	(cd cmd/flarectl && go get && go build -o ../../bin/flarectl)
	(cd cmd/flareapi && go get && go build -o ../../bin/flareapi)
	(cd cmd/flared && go get && go build -o ../../bin/flared)

test:
	go get github.com/stretchr/testify
	go test -v ./pkg/...

unit_tests:
	tests/setup.sh
	make tests

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog
