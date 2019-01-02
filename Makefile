# Copyright 2016 Francisco Souza. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: test

test: lint gotest

coverage: testdeps
	go test -race -coverprofile coverage.txt -covermode=atomic

lint: testdeps
	golangci-lint run --enable-all -D errcheck -D lll -D dupl -D gochecknoglobals -D scopelint --deadline 5m

gotest: testdeps
	go test -race

testdeps:
	go mod download
	cd /tmp && go get github.com/golangci/golangci-lint/cmd/golangci-lint
