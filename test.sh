#!/bin/bash
set -x
./build.sh

export KOPI_LOG_LEVEL=DEBUG
export KOPI_PASSWORD=password
robot --debugfile debug.log \
    test/indexing.robot \
    test/diffing.robot \
    test/store.robot \
    test/restore.robot \
    test/manifest.robot
