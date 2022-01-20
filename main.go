// Copyright (C) 2020 Alessandro Segala (ItalyPaleAle)
// License: MIT

package main

// Import the package to access the Wasm environment
import (
	"io/ioutil"
	"net/http"
	"syscall/js"
)

// async function MyFunc() {
// 	try {
// 			const response = await MyGoFunc('https://api.taylor.rest/')
// 			const message = await response.json()
// 			console.log(message)
// 	} catch (err) {
// 			console.error('Caught exception', err)
// 	}
// }

// Main function: it sets up our Wasm application
func main() {
	// Define the function "MyGoFunc" in the JavaScript scope
	js.Global().Set("MyGoFunc", MyGoFunc())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}

// MyGoFunc fetches an external resource by making a HTTP request from Go
// The JavaScript method accepts one argument, which is the URL to request
func MyGoFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		requestUrl := args[0].String()
		// We need to return a Promise because HTTP requests are blocking in Go
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]
			go func() {
				res, err := http.DefaultClient.Get(requestUrl)
				if err != nil {
					// Handle errors: reject the Promise if we have an error
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}
				defer res.Body.Close()

				// Read the response body
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					// Handle errors here too
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}

				// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
				arrayConstructor := js.Global().Get("Uint8Array")
				dataJS := arrayConstructor.New(len(data))
				js.CopyBytesToJS(dataJS, data)

				// Create a Response object and pass the data
				responseConstructor := js.Global().Get("Response")
				response := responseConstructor.New(dataJS)

				// Resolve the Promise
				resolve.Invoke(response)
			}()
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
