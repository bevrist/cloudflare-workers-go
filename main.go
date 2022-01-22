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
	w.Header().Set("testtttt", "passssss")
	w.Header().Add("testtttt", "22222")
	w.Header().Add("testtttt", "passsss2")
	resp, _ := http.Get("https://random-data-api.com/api/stripe/random_stripe")
	reqBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, string(reqBody))
}

// Main function: it sets up our Wasm application
func main() {
	// Define the function "MyGoFunc" in the JavaScript scope
	js.Global().Set("WorkerWrapper", WorkerHandlerWrapper())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}

// MyGoFunc fetches an external resource by making a HTTP request from Go
// The JavaScript method accepts one argument, which is the URL to request
func WorkerHandlerWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// requestUrl := args[0].String()
		// We need to return a Promise because HTTP requests are blocking in Go
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			// reject := args[1]  // can reject the promise with this followed by a return:
			// 		reject.Invoke(js.Global().Get("Error").New(err.Error()))
			go func() {
				var w http.ResponseWriter = new(jsResponseWriter)
				var r http.Request
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
