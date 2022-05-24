all: help

GO_BUILD_ARGS=CGO_ENABLED=0 GO111MODULE=on

.PHONY: help
help:
	@echo -e "help:\tYou're reading it!"
	@echo -e "build:\tBuild the codebase"
	@echo -e "clean:\tClean up build artifacts"

.PHONY: build
build:
	$(GO_BUILD_ARGS) go build -o jiraCommentsReport

.PHONY: clean
clean:
	$(GO_BUILD_ARGS) go clean
