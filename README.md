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
- convert incoming request to http.Request
- convert outgoing http.ResponseWriter to javascript string
