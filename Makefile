.PHONY: build clean

GIT_ROOT := $(shell git rev-parse --show-toplevel)
BINAREY_NAME := tsdw

build:
	go build -ldflags="-s -w" -o $(GIT_ROOT)/$(BINAREY_NAME) $(GIT_ROOT)/main.go

clean:
	@-rm -f $(BINAREY_NAME)