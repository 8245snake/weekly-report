set GOOS=js
set GOARCH=wasm
go build  -tags=wasm -o main.wasm
pause