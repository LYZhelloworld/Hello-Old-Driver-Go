@echo off
set GO111MODULE=on
go test -v llss/internal/scanner llss/internal/analyzer llss/internal/domainscanner
