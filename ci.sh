#!/bin/bash
set -ex
sudo apt-get update
sudo apt-get -y install python3 python3-pip
sudo -H pip3 install robotframework
python --version
go get ./...
./test.sh
