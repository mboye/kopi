#!/bin/bash
set -x
./build.sh
robot \
    test/indexing.txt \
    test/diffing.txt
