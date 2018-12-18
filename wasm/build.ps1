$env:GOARCH="wasm";
$env:GOOS="js";
go build -o public/index.wasm wasm/main.go wasm/chess.go;
$env:GOARCH="";
$env:GOOS="";
