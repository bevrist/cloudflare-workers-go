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
  const response = await WorkerWrapper()
  console.log(response.response)
  console.log(response.response.headers)
  // console.log(JSON.parse(response).response)
  return new Response(response.body, response.response);
}
