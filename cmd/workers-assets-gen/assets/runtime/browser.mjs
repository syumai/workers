export async function loadModule() {
  return await WebAssembly.compileStreaming(fetch("./build/app.wasm"));
}

export function createRuntimeContext({ binding }) {
  return {
    binding,
  };
}
