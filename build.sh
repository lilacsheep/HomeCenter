#!/usr/bin/env bash

rm -rf build && mkdir -p build
rm -rf packed/data.go
# 依赖 goframe 框架工具 进行静态文件注入
gf pack public packed/data.go -n=packed

GOOS=linux GOARCH=amd64 go build -o build/proxy
GOOS=windows GOARCH=amd64 go build -o build/proxy.exe
GOOS=linux GOARCH=arm go build -o build/proxy_arm_v7
GOOS=linux GOARCH=arm64 go build -o build/proxy_arm_v8