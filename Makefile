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
	@bazel run //cmd
	# @bazel run $(shell sudo sh "$(PWD)"/env.sh 1) //cmd

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
	@bazel test --test_output=summary --test_timeout=2 -t-  //...
	# @bazel test $(shell sudo sh "$(PWD)"/env.sh 2) --test_output=summary --test_timeout=2 -t-  //...

## setup
.PHONY:
setup:
	@docker run --name mongo4 -d -p 27017:27017 --rm mongo