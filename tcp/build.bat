@echo off
set GOARCH=386
go build -o out/client.exe client/main.go
go build -o out/server.exe server/main.go
