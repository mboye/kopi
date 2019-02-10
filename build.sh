#!/bin/bash
set -ex
mkdir -p bin
go build -o bin/kopi-index cmd/index/index.go
go build -o bin/kopi-diff cmd/diff/diff.go
go build -o bin/kopi-store cmd/store/store.go
go build -o bin/kopi-restore cmd/restore/restore.go
go build -o bin/kopi-manifest cmd/manifest/manifest.go
