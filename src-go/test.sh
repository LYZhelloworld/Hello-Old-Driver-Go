#/bin/sh
export GOPATH=$PWD
go test -v scanner analyzer domain_scanner
