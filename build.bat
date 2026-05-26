set GOOS=windows
set GOARCH=arm64
go build -o dist/win_64.exe ./src

set GOARCH=386
go build -o dist/win_32.exe ./src

set GOOS=linux
set GOARCH=amd64
go build -o dist/linux.bin ./src

set GOOS=darwin 
set GOARCH=arm64
go build -o dist/mac-arm.bin ./src