#/qompassai/adns/Makefile
# -----------------------
# Copyright (C) 2025 Qompass AI, All rights reserved
build:
	go build
	go vet ./...

test:
	go test -race -shuffle=on -coverprofile cover.out -covermode atomic
	go tool cover -html=cover.out -o cover.html

check:
	GOARCH=386 go vet
	staticcheck -checks inherit,-SA1019 ./...

check-shadow:
	go vet -vettool=$$(which shadow) ./... 2>&1 | grep -v '"err"'

buildall:
	GOOS=linux GOARCH=arm go build
	GOOS=linux GOARCH=arm64 go build
	GOOS=linux GOARCH=amd64 go build
	GOOS=linux GOARCH=386 go build
	GOOS=openbsd GOARCH=amd64 go build
	GOOS=freebsd GOARCH=amd64 go build
	GOOS=netbsd GOARCH=amd64 go build
	GOOS=darwin GOARCH=amd64 go build
	GOOS=dragonfly GOARCH=amd64 go build
	GOOS=illumos GOARCH=amd64 go build
	GOOS=solaris GOARCH=amd64 go build
	GOOS=aix GOARCH=ppc64 go build
	GOOS=windows GOARCH=amd64 go build

fmt:
	gofmt -w -s *.go */*/*.go
