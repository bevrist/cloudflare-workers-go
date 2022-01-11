import './polyfill';
import './wasm_exec';

addEventListener("fetch", (event) => {
  event.respondWith(handleRequest(event.request));
});

async function handleRequest(request) {
  const go = new Go();
  const instance = await WebAssembly.instantiate(WASM, go.importObject);
  go.run(instance);
  go.run(instance);
  return new Response("Make sure a and b are numbers\n", { status: 200 });
}
