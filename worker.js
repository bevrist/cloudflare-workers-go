import './polyfill';
import './wasm_exec';

addEventListener("fetch", (event) => {
  event.respondWith(handleRequest(event.request));
});

function test() {
  console.log("test1 succeeded");
}

var lol = "hoo"
function test2() {
  return "a test2"
}

async function handleRequest(request) {
  // provide a reference for all js objects so they are not tree-shaken out of the final build
  var antiTreeShake = { test, test2, lol }

  const go = new Go();
  const instance = await WebAssembly.instantiate(WASM, go.importObject);
  go.run(instance);
  const body = await request.text()
  //TODO: try to use js arrayBuffer and golang CopyBytesToGo instead of text()
  console.log(body)
  var headerKeys = []
  var headerValues = []
  for (var pair of request.headers.entries()) {
    headerKeys.push(pair[0])
    headerValues.push(pair[1])
  }
  // console.log(JSON.stringify(headerKeys))
  // console.log(JSON.stringify(headerValues))
  // pass request objects to golang func
  // all parameters are contained in args[] as a js.Value
  const response = await WorkerHandlerWrapper(request, body, headerKeys.length, headerKeys, headerValues)
  // console.log(Object.fromEntries(request.headers))
  return new Response(response.body, response.response);
}


// // Create a test Headers object
// var myHeaders = new Headers();
// myHeaders.append('Content-Type', 'text/xml');
// myHeaders.append('Vary', 'Accept-Language');

// // Display the key/value pairs
// for (var pair of myHeaders.entries()) {
//    console.log(pair[0]+ ': '+ pair[1]);
// }