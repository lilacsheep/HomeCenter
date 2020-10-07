@echo off
SET GOOS=linux
SET GOARCH=amd64
go build -o build/proxy


SET GOOS=windows
SET GOARCH=amd64
go build -o build/proxy.exe


SET GOOS=linux
SET GOARCH=arm
go build -o build/proxy_arm_v7

SET GOOS=linux
SET GOARCH=arm64
go build -o build/proxy_arm_v8