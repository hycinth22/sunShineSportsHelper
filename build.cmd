@echo off

rem win64
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o build/jkwx.exe

xcopy /-Y /R build_assets build
