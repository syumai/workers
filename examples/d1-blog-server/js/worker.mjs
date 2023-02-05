import "./polyfill_performance.js";
import "./wasm_exec.js";
// import "./wasm_exec_go.js";
import mod from "../dist/app.wasm";

const go = new Go();

const readyPromise = new Promise((resolve) => {
  globalThis.ready = resolve;
});

const load = WebAssembly.instantiate(mod, go.importObject).then((instance) => {
  go.run(instance);
  return instance;
});

async function processRequest(req) {
  await load;
  await readyPromise;
  return handleRequest(req);
}

export default {
  async fetch(req, env) {
    for (const [key, value] of Object.entries(env)) {
      globalThis[key] = value;
    }
    return processRequest(req);
  }
}
