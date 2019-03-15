#!/bin/bash
set -ex

build() {
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags -static' $@
}

output_dir=bin/linux
mkdir -p "$output_dir"
build -o $output_dir/kopi-index cmd/index/index.go
build -o $output_dir/kopi-diff cmd/diff/diff.go
build -o $output_dir/kopi-store cmd/store/store.go
build -o $output_dir/kopi-restore cmd/restore/restore.go
build -o $output_dir/kopi-manifest cmd/manifest/manifest.go
