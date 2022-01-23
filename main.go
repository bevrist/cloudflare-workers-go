// Copyright (C) 2020 Alessandro Segala (ItalyPaleAle)
// License: MIT

package main

// Import the package to access the Wasm environment
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall/js"
)

type jsResponseWriter struct {
	http.ResponseWriter
	headers    http.Header
	body       []byte
	statusCode int
}

func (w *jsResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return 200, nil
}

func (w *jsResponseWriter) Header() http.Header {
	w.headers = make(http.Header)
	return w.headers
}

func (w *jsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func WorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("test", "pass")
	resp, _ := http.Get("https://random-data-api.com/api/stripe/random_stripe")
	reqBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, string(reqBody))
}

// Main function: it sets up our Wasm application
func main() {
	// Define the function "WorkerHandlerWrapper" in the JavaScript scope
	js.Global().Set("WorkerHandlerWrapper", WorkerHandlerWrapper())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}

// MyGoFunc fetches an external resource by making a HTTP request from Go
// The JavaScript method accepts one argument, which is the URL to request
func WorkerHandlerWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// get passed parameters and convert to golang types
		request := args[0] // the request object was 1st parameter
		// body := args[1].String() // the body string was the 2nd parameter

		// extract request headers
		headerLen := args[2].Int() // the 3-5 parameters are headers
		headerKeys := args[3]
		headerValues := args[4]
		var reqHeaders http.Header = make(http.Header)
		for i := 0; i < headerLen; i++ {
			reqHeaders[headerKeys.Index(i).String()] = []string{headerValues.Index(i).String()}
		}

		// create golang request
		var r http.Request
		r.Method = request.Get("method").String()
		r.Header = reqHeaders
		// make URL object
		// scheme = regex: `(http.?:\/\/)`
		// host = regex: `([\w-]+(?:\.[\w-]+)*(?::\d+)?)` //includes potential :port
		// path = regex: `(\/.*)\?` //this leaves out ? parameters (needed?)
		//		if path re is empty its relative address, set raw path string as path
		// use URL.SetPath(path) to set path and RawPath

		// r.URL = request.Get("url").String()
		// url.URL =

		// 	type URL struct {
		// 		Scheme      string
		//// 		Opaque      string    // encoded opaque data
		//// 		User        *Userinfo // username and password information
		// 		Host        string    // host or host:port
		/// 		Path        string    // path (relative paths may omit leading slash)
		/// 		RawPath     string    // encoded path hint (see EscapedPath method)
		//// 		ForceQuery  bool      // append a query ('?') even if RawQuery is empty
		//// 		RawQuery    string    // encoded query values, without '?'
		//// 		Fragment    string    // fragment for references, without '#'
		//// 		RawFragment string    // encoded fragment hint (see EscapedFragment method)
		// }

		// We need to return a Promise because HTTP requests are blocking in Go
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			// reject := args[1]  // can reject the promise with this followed by a return:
			// 		reject.Invoke(js.Global().Get("Error").New(err.Error()))
			go func() {
				var w http.ResponseWriter = new(jsResponseWriter)
				WorkerHandler(w, &r)
				a := w.(*jsResponseWriter)
				bodyInit := make(map[string]interface{})
				bodyInit["body"] = string(a.body)
				responseInit := make(map[string]interface{})
				if a.statusCode == 0 {
					responseInit["status"] = 200
				} else {
					responseInit["status"] = a.statusCode
				}
				//get headers from golang
				headers := make(map[string]interface{})
				for key := range a.headers {
					headers[key] = strings.Join(a.headers.Values(key), ", ")
				}
				responseInit["headers"] = headers
				bodyInit["response"] = responseInit
				resolve.Invoke(js.ValueOf(bodyInit))
			}()
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
