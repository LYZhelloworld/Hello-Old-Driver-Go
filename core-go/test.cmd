@echo off
set GOPATH=%cd%
go test -v scanner analyzer domain_scanner
