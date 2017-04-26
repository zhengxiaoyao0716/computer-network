@echo off
set GOARCH=386
go build -o out/tcp-client.exe tcp-client/main.go
go build -o out/tcp-server.exe tcp-server/main.go
