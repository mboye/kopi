#!/bin/bash
set -ex
sudo apt-get -y install python3
pip install robotframework
go get ./...
./test.sh
