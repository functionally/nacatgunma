#!/usr/bin/env nix-shell
#!nix-shell -i bash -p go golangci-lint

go build -o nacatgunma main.go
