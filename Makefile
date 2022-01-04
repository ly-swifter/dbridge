SHELL=/usr/bin/env bash

.PHONY: clean
clean:
	rm -f dbridge

.PHONY: all
all:
	go build -o dbridge ./cmd/dbridge