#!/bin/bash
set -ex
sudo apt-get update &>/dev/null
sudo apt-get -y install python3 python3-pip &>/dev/null
sudo -H pip3 install robotframework
python --version
go get ./...
./test.sh
