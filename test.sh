#!/bin/bash
set -x
./build.sh
robot --debugfile debug.log \
    test/indexing.robot \
    test/diffing.robot \
    test/store.robot \
    test/restore.robot
