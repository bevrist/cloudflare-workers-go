# Cloudflare-Workers-Go

use golang to power your cloudflare workers.

## Build:
`GOARCH=wasm GOOS=js go build -o main.wasm main.go && npm run build && wrangler dev`
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
- find better fix for golang wasm compilation
- convert project to library
- extract body and headers for js response
- implement shim for cloudflare features
  - workers kv
  - environment vars
  - durable objects
