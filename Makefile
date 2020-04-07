all: build

GO_BUILD_ARGS=CGO_ENABLED=0 GO111MODULE=on

.PHONY: build
build:
	$(GO_BUILD_ARGS) go build -o jiraCommentsReport
