#!/usr/bin/env nix-shell
#!nix-shell -i "make -f" -p go golint

PKGS=./cardano ./cmd ./header ./ipfs ./key ./ledger ./rdf ./tgdh

SRCS=$(shell find $(PKGS) -type f -name \*.go)

nacatgunma: main.go $(SRCS)
	GOPATH= go build -o nacatgunma $<

test: $(SRCS)
	GOPATH= go test -v $(PKGS)

format: $(SRCS)
	GOPATH= go fmt $(PKGS)

.SUFFIXES:

.PHONY: test format