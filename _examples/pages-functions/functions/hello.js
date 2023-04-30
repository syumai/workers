import "../../assets/polyfill_performance.js";
import "../../assets/wasm_exec.js";
import mod from "../build/app.wasm";

const go = new Go();

const readyPromise = new Promise((resolve) => {
  globalThis.ready = resolve;
});

const load = WebAssembly.instantiate(mod, go.importObject).then((instance) => {
  go.run(instance);
  return instance;
});

export const onRequest = async (ctx) => {
  await load;
  await readyPromise;
  const {
    request,
    env,
  } = ctx;
  return handleRequest(request, { env, ctx });
}