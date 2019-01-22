#!/bin/bash
set -x
./build.sh
robot \
    test/indexing.robot \
    test/diffing.robot \
    test/store.robot \
    test/restore.robot
