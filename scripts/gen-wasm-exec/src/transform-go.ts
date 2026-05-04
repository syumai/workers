import { transformWasmExec } from "./transform-wasm-exec.ts";

const GLOBAL_PROXY_DECL = `const globalProxy = new Proxy(globalThis, {
\tget(target, prop) {
\t\tif (prop === 'context') {
\t\t\treturn context;
\t\t}
\t\treturn Reflect.get(target, prop, target);
\t}
})`;

export function transformGo(rawSource: string): string {
	return transformWasmExec(rawSource, {
		flavor: "go",
		globalName: "globalThis",
		expectedStandaloneGlobals: 2,
		globalProxyDecl: GLOBAL_PROXY_DECL,
	});
}
