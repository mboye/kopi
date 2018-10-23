#!/bin/bash
set -x
./build.sh
pybot \
    test/indexing.txt \
    test/diffing.txt
