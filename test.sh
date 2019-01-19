#!/bin/bash
set -x
./build.sh
robot \
    test/indexing.robot \
    test/diffing.robot
