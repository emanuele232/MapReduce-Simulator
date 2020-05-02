#!/bin/zsh

env GOOS=darwin GOARCH=amd64 go build -o osx64Build main.go
env GOOS=linux GOARCH=amd64 go build -o linux64Build main.go
env GOOS=windows GOARCH=amd64 go build -o win64Build.exe main.go
chmod +x osx64Build
chmod +x win64Build.exe
