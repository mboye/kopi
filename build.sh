#!/bin/bash
set -ex
mkdir -p bin
go build -o bin/kopi-index cmd/index/index.go
