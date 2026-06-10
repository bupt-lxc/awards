#!/bin/sh
set -e

export GOPROXY=https://goproxy.cn,direct
go mod download
go mod tidy
go build -o backend .
exec ./backend
