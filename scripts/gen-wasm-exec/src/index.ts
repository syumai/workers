import { writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { fetchText, goWasmExecURL, tinygoWasmExecURL } from "./fetch.ts";
import { transformGo } from "./transform-go.ts";
import { transformTinygo } from "./transform-tinygo.ts";

interface Args {
	go?: string;
	tinygo?: string;
	out?: string;
}

function parseArgs(argv: string[]): Args {
	const args: Args = {};
	for (let i = 0; i < argv.length; i++) {
		const a = argv[i];
		const next = argv[i + 1];
		switch (a) {
			case "--go":
				if (!next) throw new Error("--go requires a value");
				args.go = next;
				i++;
				break;
			case "--tinygo":
				if (!next) throw new Error("--tinygo requires a value");
				args.tinygo = next;
				i++;
				break;
			case "--out":
				if (!next) throw new Error("--out requires a value");
				args.out = next;
				i++;
				break;
			default:
				throw new Error(`unknown argument: ${a}`);
		}
	}
	return args;
}

async function main() {
	const args = parseArgs(process.argv.slice(2));
	if (!args.go) throw new Error("--go <version> is required (e.g. --go 1.26.2)");
	if (!args.tinygo)
		throw new Error("--tinygo <version> is required (e.g. --tinygo 0.41.1)");

	const here = path.dirname(fileURLToPath(import.meta.url));
	// here is scripts/gen-wasm-exec/src ; default output is repo's cmd/workers-assets-gen/assets
	const outDir = args.out
		? path.resolve(args.out)
		: path.resolve(here, "../../../cmd/workers-assets-gen/assets");

	const goURL = goWasmExecURL(args.go);
	const tinygoURL = tinygoWasmExecURL(args.tinygo);

	console.log(`fetching ${goURL}`);
	const goSrc = await fetchText(goURL);
	console.log(`fetching ${tinygoURL}`);
	const tinygoSrc = await fetchText(tinygoURL);

	console.log("patching Go wasm_exec.js");
	const goOut = transformGo(goSrc);
	console.log("patching TinyGo wasm_exec.js");
	const tinygoOut = transformTinygo(tinygoSrc);

	const goPath = path.join(outDir, "wasm_exec_go.js");
	const tinygoPath = path.join(outDir, "wasm_exec_tinygo.js");
	await writeFile(goPath, goOut);
	await writeFile(tinygoPath, tinygoOut);
	console.log(`wrote ${goPath}`);
	console.log(`wrote ${tinygoPath}`);
}

main().catch((err) => {
	console.error(err);
	process.exit(1);
});
