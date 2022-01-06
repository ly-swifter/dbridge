SHELL=/usr/bin/env bash

GOCC?=go

.PHONY: clean
clean:
	rm -f lorry

.PHONY: all
all:
	$(GOCC) build -o lorry ./cmd/dbridge

api-gen:
	$(GOCC) run ./gen/api
	goimports -w api
	goimports -w api
.PHONY: api-gen