set GOOS=windows
set GOARCH=arm64
go build -o dist/win_64.exe ./src

set GOARCH=386
go build -o dist/win_32.exe ./src