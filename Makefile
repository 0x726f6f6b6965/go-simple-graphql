PROJECTNAME := $(shell basename "$(PWD)")

#set our own default goal
.PHONY: default
default: help

.PHONY: help
help: Makefile
	@echo
	@echo "Choose a command to run in $(PROJECTNAME):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo

## gen-mod: generate graphQL model
.PHONY: gen-mod
gen-mod:
	@go generate ./graph/resolver.go

## run: run the graphQL server
.PHONY: run
run:
	@go run server.go

## clean
.PHONY: clean
clean:
	@bazel clean

## build
.PHONY: build
build:
	@bazel run //:gazelle-update-repos 
	@bazel run //:gazelle

## test
.PHONY: test
test:
	@bazel test --test_output=summary --test_timeout=2 -t- //...