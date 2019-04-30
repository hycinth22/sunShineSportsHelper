@echo off

rem win64
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o releases/sunShineSportsHelper-win-amd64.exe

rem linux 64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o releases/sunShineSportsHelper-linux-amd64

xcopy /-Y /R build_assets assets\releases\