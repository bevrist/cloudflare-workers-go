# Cloudflare-Workers-Go

use golang to power your cloudflare workers.

## Build:
### golang:
```bash
cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./
npm install
npm run build
GOARCH=wasm GOOS=js go build -o main.wasm main.go
```

### tinygo:
```bash
# TODO finish this
# cp tinigo wasm_exec.js
npm install
npm run build
# tinygo build ...
```

## TODO:
- make shim library for `net/http`
- use go:build to filter shim or real `net/http` based on environment https://github.com/golang/go/blob/master/src/syscall/js/js.go 
- make cloudflare workers js lib
