SHELL=/usr/bin/env bash

.PHONY: clean
clean:
	rm -f lorry

.PHONY: all
all:
	go build -o lorry ./cmd/dbridge