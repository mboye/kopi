#!/bin/bash
set -x
./build.sh
pybot test/indexing.txt
