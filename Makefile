.DEFAULT_GOAL := help
CURRENTTAG:=$(shell git describe --tags --abbrev=0)
NEWTAG ?= $(shell bash -c 'read -p "Please provide a new tag (currnet tag - ${CURRENTTAG}): " newtag; echo $$newtag')
GOFLAGS=-mod=mod

#help: @ List available tasks
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-6s\033[0m - %s\n", $$1, $$2}'

#clean: @ Cleanup
clean:
	@rm -f ./go-http-instrumented-roundtripper

#get: @ Download and install dependency packages
get:
	@export GOFLAGS=$(GOFLAGS); go get . ; go mod tidy

#update: @ Update dependencies to latest versions
update:
	@go get -u; go mod tidy

#test: @ Run tests
test:
	@export GOPRIVATE=$(GOPRIVATE); go generate
	@export GOPRIVATE=$(GOPRIVATE); export GOFLAGS=$(GOFLAGS); go test $(go list ./... | grep -v /internal/setup)

#build: @ Build binary
build:
	@export GOPRIVATE=$(GOPRIVATE); go generate
	@export GOPRIVATE=$(GOPRIVATE); export GOFLAGS=$(GOFLAGS); export CGO_ENABLED=0; go build -a -o main main.go

release: build
	$(eval NT=$(NEWTAG))
	@echo -n "Are you sure to create and push ${NT} tag? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo ${NT} > ./version.txt
	@git add -A
	@git commit -a -s -m "Cut ${NT} release"
	@git tag -a -m "Cut ${NT} release" ${NT}
	@git push origin ${NT}
	@git push
	@echo "Done."

version:
	@echo $(shell git describe --tags --abbrev=0)
