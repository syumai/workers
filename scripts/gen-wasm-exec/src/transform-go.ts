import { transformWasmExec } from "./transform-wasm-exec.ts";

export function transformGo(rawSource: string): string {
	return transformWasmExec(rawSource, {
		flavor: "go",
		globalName: "globalThis",
		expectedStandaloneGlobals: 2,
	});
}
