// Copyright (C) 2020 Alessandro Segala (ItalyPaleAle)
// License: MIT

package main

// Import the package to access the Wasm environment
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"
)

type jsResponseWriter struct {
	http.ResponseWriter
	header     http.Header
	body       []byte
	statusCode int
}

func (w *jsResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return 200, nil
}

func (w *jsResponseWriter) Header() http.Header {
	w.header = make(http.Header)
	return w.header
}

func (w *jsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func WorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("testtttt", "passssss")
	resp, _ := http.Get("https://random-data-api.com/api/stripe/random_stripe")
	reqBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "hi: "+string(reqBody))
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
			// reject := args[1]
			go func() {
				// res, err := http.DefaultClient.Get("https://random-data-api.com/api/stripe/random_stripe")
				// if err != nil {
				// 	// Handle errors: reject the Promise if we have an error
				// 	errorConstructor := js.Global().Get("Error")
				// 	errorObject := errorConstructor.New(err.Error())
				// 	reject.Invoke(errorObject)
				// 	return
				// }
				// defer res.Body.Close()

				// // Read the response body
				// data, err := ioutil.ReadAll(res.Body)
				// if err != nil {
				// 	// Handle errors here too
				// 	errorConstructor := js.Global().Get("Error")
				// 	errorObject := errorConstructor.New(err.Error())
				// 	reject.Invoke(errorObject)
				// 	return
				// }

				var w http.ResponseWriter = new(jsResponseWriter)
				var r http.Request
				WorkerHandler(w, &r)
				a := w.(*jsResponseWriter)
				// resolve.Invoke(string(a.body))
				respond := make(map[string]interface{})
				respond["body"] = string(a.body)
				respond["response"] = "a=b"
				resolve.Invoke(js.ValueOf(respond))
			}()
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
