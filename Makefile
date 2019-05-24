# Task
#===============================================================

fmt:
	go fmt $$(go list ./... | tr '\n' ' ')

test:
	go test $$(go list ./... | tr '\n' ' ')

build:
	go build -o ./bin/kangol *.go

update:
	go get -u=patch

update_test: update test

.PHONY: fmt test build update