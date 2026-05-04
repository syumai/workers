const RAW_GITHUB = "https://raw.githubusercontent.com/";

export function goWasmExecURL(version: string): string {
	return new URL(
		`golang/go/refs/tags/go${version}/lib/wasm/wasm_exec.js`,
		RAW_GITHUB,
	).toString();
}

export function tinygoWasmExecURL(version: string): string {
	return new URL(
		`tinygo-org/tinygo/refs/tags/v${version}/targets/wasm_exec.js`,
		RAW_GITHUB,
	).toString();
}

export async function fetchText(url: string): Promise<string> {
	const res = await fetch(url);
	if (!res.ok) {
		throw new Error(`failed to fetch ${url}: ${res.status} ${res.statusText}`);
	}
	return await res.text();
}
